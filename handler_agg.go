package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Marcos-Pablo/go-blog-aggregator/internal/rss"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])

	if err != nil {
		return fmt.Errorf("Invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s...\n", timeBetweenReqs.String())
	fmt.Println()

	ticker := time.NewTicker(timeBetweenReqs)
	defer ticker.Stop()

	for range ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func scrapeFeeds(s *state) {
	feed, err := s.queries.GetNextFeedToFetch(context.Background())

	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	var errMarkFetched error
	var errFetchFeed error
	var rssFeed *rss.RSSFeed

	go func() {
		_, errMarkFetched = s.queries.MarkFeedFetched(context.Background(), feed.ID)
		wg.Done()
	}()

	go func() {
		rssFeed, errFetchFeed = rss.FetchFeed(context.Background(), feed.Url)
		wg.Done()
	}()

	wg.Wait()

	if errMarkFetched != nil || errFetchFeed != nil {
		log.Printf("Couldn't scrape feed %s: %v", feed.Name, errors.Join(errMarkFetched, errFetchFeed))
		return
	}

	fmt.Println("RSS feed fetched successfully:")
	for _, item := range rssFeed.Channel.Item {
		if item.Title == "" {
			continue
		}
		fmt.Printf("* Title:      %s\n", item.Title)
	}
	fmt.Println()
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
	fmt.Println("=====================================")
}
