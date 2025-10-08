package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JoeVinten/blogregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])

	if err != nil {
		return fmt.Errorf("invalidation error, %w", err)
	}

	fmt.Printf("Collecting feeds every %s...", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {

	feed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		fmt.Println("could not get next feed to fetch", err)
		return err
	}

	log.Println("Found a feed to fetch!")

	scrapeFeed(s.db, feed)

	return nil

}

func scrapeFeed(db *database.Queries, feed database.Feed) {

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)

	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)

	if err != nil {
		log.Printf("Couldn't fetch feed from url %s: %v", feed.Url, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))

}
