package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, _ command) error {
	err := s.queries.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Reset executed succesfully")
	return nil
}
