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
		res, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
		user = res
		userErr = err
		wg.Done()
	}()

	go func() {
		res, err := s.db.GetFeedByUrl(ctx, url)
		feed = res
		feedErr = err
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

	follow, err := s.db.CreateFeedFollow(context.Background(), params)

	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	fmt.Println("Successfully followed feed")
	fmt.Printf("* ID:        %s\n", follow.ID)
	fmt.Printf("* feed:      %s\n", feed.Name)
	fmt.Printf("* user:      %s\n", user.Name)

	return nil
}
