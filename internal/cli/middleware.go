package cli

import (
	"bootdevBlogAggerator/internal/database"
	"context"
	"fmt"
)

func middlewareLoggedIn(handler func(s *State, cmd Command, currentUser database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return fmt.Errorf("Something went wrong fetching current user: %v\n", err)
		}
		return handler(s, cmd, currentUser)
	}
}
