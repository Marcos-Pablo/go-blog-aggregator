package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Marcos-Pablo/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}

	name, url := cmd.args[0], cmd.args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)

	if err != nil {
		return err
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)

	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	fmt.Println("Feed created successfully")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	printBatchFeed(feeds)
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

func printBatchFeed(feeds []database.GetFeedsRow) {
	for _, feed := range feeds {
		fmt.Printf("* ID:            %s\n", feed.ID)
		fmt.Printf("* Created:       %v\n", feed.CreatedAt)
		fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
		fmt.Printf("* Name:          %s\n", feed.Name)
		fmt.Printf("* URL:           %s\n", feed.Url)
		fmt.Printf("* UserID:        %s\n", feed.UserID)
		fmt.Printf("* UserName:      %s\n", feed.UserName)
		fmt.Println()
		fmt.Println("=====================================")
	}
}
