package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Marcos-Pablo/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	url := cmd.args[0]
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)
	var user database.User
	var feed database.Feed
	var userErr error
	var feedErr error

	go func() {
		res, err := s.queries.GetUser(ctx, s.cfg.CurrentUserName)
		user = res
		userErr = err
		if err != nil {
			cancel()
		}
		wg.Done()
	}()

	go func() {
		res, err := s.queries.GetFeedByUrl(ctx, url)
		feed = res
		feedErr = err
		if err != nil {
			cancel()
		}
		wg.Done()
	}()

	wg.Wait()

	if userErr != nil || feedErr != nil {
		return fmt.Errorf("couldn't follow feed: %w", errors.Join(userErr, feedErr))
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follow, err := s.queries.CreateFeedFollow(ctx, params)

	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	fmt.Println("Successfully followed feed")
	fmt.Printf("* ID:        %s\n", follow.ID)
	fmt.Printf("* feed:      %s\n", feed.Name)
	fmt.Printf("* user:      %s\n", user.Name)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	user, err := s.queries.GetUser(context.Background(), s.cfg.CurrentUserName)

	if err != nil {
		return fmt.Errorf("couldn't get the followed feeds for the current user: %w", err)
	}

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
