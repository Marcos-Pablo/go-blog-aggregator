package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Marcos-Pablo/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	url := cmd.args[0]

	feed, err := s.queries.GetFeedByUrl(context.Background(), url)

	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follow, err := s.queries.CreateFeedFollow(context.Background(), params)

	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	fmt.Println("Successfully followed feed")
	fmt.Printf("* ID:        %s\n", follow.ID)
	fmt.Printf("* feed:      %s\n", feed.Name)
	fmt.Printf("* user:      %s\n", user.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.queries.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("couldn't get the followed feeds for the current user: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Println("Feeds you're currently following:")

	for _, follow := range follows {
		fmt.Printf("*      %s\n", follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	url := cmd.args[0]

	feed, err := s.queries.GetFeedByUrl(context.Background(), url)

	if err != nil {
		return fmt.Errorf("couldn't unfollow feed: %w", err)
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.queries.DeleteFeedFollow(context.Background(), params)

	if err != nil {
		return fmt.Errorf("couldn't unfollow feed: %w", err)
	}

	fmt.Printf("%s unfollowed successfully!\n", feed.Name)

	return nil
}
