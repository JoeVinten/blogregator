package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JoeVinten/blogregator/internal/database"
	"github.com/google/uuid"
)

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
