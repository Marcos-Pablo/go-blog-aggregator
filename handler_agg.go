package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Marcos-Pablo/go-blog-aggregator/internal/database"
	"github.com/Marcos-Pablo/go-blog-aggregator/internal/rss"
	"github.com/google/uuid"
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
	// TODO: batch insert
	for _, item := range rssFeed.Channel.Item {
		if item.Title == "" {
			continue
		}

		pubDate, err := parseDate(item.PubDate, []string{time.RFC3339, time.RFC1123Z})

		if err != nil {
			log.Println(err)
			continue
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		}

		_, err = s.queries.CreatePost(context.Background(), params)

		if err != nil {
			log.Printf("couldn't create post: %v", err)
			continue
		}

		fmt.Println("Post created successfully")
	}
}

func parseDate(strDate string, layouts []string) (time.Time, error) {
	for _, layout := range layouts {
		date, err := time.Parse(layout, strDate)

		if err == nil {
			return date, nil
		}

	}
	return time.Time{}, fmt.Errorf("couldn't parse the date with any of the provided layouts: %s", strings.Join(layouts, ","))
}
