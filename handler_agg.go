package main

import (
	"context"
	"fmt"

	"github.com/Marcos-Pablo/go-blog-aggregator/internal/rss"
)

func handlerAgregate(s *state, cmd command) error {
	// if len(cmd.args) != 1 {
	// 	return fmt.Errorf("usage: %s <url>", cmd.name)
	// }
	//
	// url := cmd.args[0]
	url := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(context.Background(), url)

	if err != nil {
		return fmt.Errorf("Couldn't fetch feed: %w", err)
	}

	printFeed(feed)
	return nil
}

func printFeed(feed *rss.RSSFeed) {
	fmt.Println("Feed:")
	fmt.Printf("- Title: %s\n", feed.Channel.Title)
	fmt.Printf("- Link: %s\n", feed.Channel.Link)
	fmt.Printf("- Description: %s\n", feed.Channel.Description)
	fmt.Println("- Items:")
	for _, item := range feed.Channel.Item {
		fmt.Printf("*  Title: %s\n", item.Title)
		fmt.Printf("*  Link: %s\n", item.Link)
		fmt.Printf("*  Descriptoin: %s\n", item.Description)
		fmt.Printf("*  PubDate: %s\n", item.PubDate)
		fmt.Println()
	}
}
