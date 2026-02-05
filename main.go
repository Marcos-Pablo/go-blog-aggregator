package main

import (
	"fmt"
	"github.com/Marcos-Pablo/go-blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error while reading config file: %v", err)
		return
	}

	username := "Marcos Pablo"
	if err = cfg.SetUser(username); err != nil {
		fmt.Printf("Error while writing config file: %v", err)
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("Error while re-reading config file: %v", err)
		return
	}

	fmt.Printf("username: %s\n", cfg.CurrentUserName)
	fmt.Printf("db_url: %s\n", cfg.DbURL)
}
