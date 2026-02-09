package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/Marcos-Pablo/go-blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: %s <limit>", cmd.name)
	}

	var limit int32 = 2
	if len(cmd.args) == 1 {
		i, err := strconv.ParseInt(cmd.args[0], 10, 32)

		if err != nil {
			log.Print("Invalid limit format, defaulting to 2")
		} else {
			limit = int32(i)
		}
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}

	posts, err := s.queries.GetPostsForUser(context.Background(), params)

	if err != nil {
		return fmt.Errorf("couldn't browse posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("0 posts found")
		return nil
	}

	fmt.Printf("Successfully retrieved %d posts\n", len(posts))
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}
