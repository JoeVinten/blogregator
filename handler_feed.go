package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/JoeVinten/blogregator/internal/database"
	"github.com/google/uuid"
)

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
