package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Marcos-Pablo/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	username := cmd.args[0]

	user, err := s.queries.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("User not found")
	} else if err != nil {
		return err
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	username := cmd.args[0]

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	}

	user, err := s.queries.CreateUser(context.Background(), params)

	if err != nil {
		return err
	}

	if err = s.cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Println("User created succesfully")
	printUser(user)

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.queries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %v\n", user.Name)
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf("-ID:         %v\n", user.ID)
	fmt.Printf("-Name:       %s\n", user.Name)
	fmt.Printf("-Created at: %v\n", user.CreatedAt)
	fmt.Printf("-Updated at: %v\n", user.UpdatedAt)
}
