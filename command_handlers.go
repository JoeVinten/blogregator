package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/JoeVinten/blogregator/internal/database"
	"github.com/google/uuid"
)

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.registeredCommands[cmd.Name]

	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("no username given")
	}

	username := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), username)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: user %s, does not exist", username)
		os.Exit(1)
	}

	err = s.cfg.SetUser(username)

	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Printf("username %s, has been set\n", username)

	return nil

}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("no username given")
	}

	username := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), username)

	if err == nil {
		fmt.Fprintf(os.Stderr, "Error: User with name '%s' already exists\n", username)
		os.Exit(1)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	})

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)

	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil

}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())

	if err != nil {
		return fmt.Errorf("failed to reset users table: %w", err)
	}

	fmt.Println("Database reset successfully!")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		return fmt.Errorf("failed to get the users: %w", err)
	}

	currentUser := s.cfg.CurrentUsername

	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("*  %s (current)\n", user.Name)
		} else {
			fmt.Printf("*  %s\n", user.Name)
		}
	}
	return nil
}

func handlerFetchFeed(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", feed)
	return nil

}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("usage: addfeed <name> <url>")
	}

	name := strings.TrimSpace(cmd.Args[0])
	url := strings.TrimSpace(cmd.Args[1])

	if name == "" || url == "" {
		return errors.New("feed name and URL cannot be empty")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to create the feed: %w", err)
	}

	_, err = s.db.CreateFeedFollowers(context.Background(), database.CreateFeedFollowersParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Feed added successfully!\n")
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("ID: %s\n", feed.ID)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		username, err := s.db.GetUsernameById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("------\n- Name: %s\n- URL: %s\n- Created by: %s\n------\n",
			feed.Name,
			feed.Url,
			username,
		)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: follow <url>")
	}

	feed, err := s.db.GetFeedsByUrl(context.Background(), cmd.Args[0])

	if err != nil {
		return fmt.Errorf("url not found in the feeds table")
	}

	feed_follow, err := s.db.CreateFeedFollowers(context.Background(), database.CreateFeedFollowersParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("%s is now following %s\n", feed_follow.UserName, feed_follow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("  * %s\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: unfollow <url>")
	}

	// Get the field ID first

	feed, err := s.db.GetFeedsByUrl(context.Background(), cmd.Args[0])

	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	return nil
}
