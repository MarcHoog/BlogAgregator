package cli

import (
	"bootdevBlogAggerator/internal/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Username is required to login '# login <username>'\n")

	} else {
		username := cmd.Args[0]

		_, err := s.db.GetUser(context.Background(), username)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Username does not exist, use '# login <username>' first\n")
		} else if err != nil {
			return fmt.Errorf("Something went wrong trying to login %s: %v\n", username, err)
		}

		err = s.cfg.SetUser(username)
		if err != nil {
			return err
		}

		println("> Logged in as: " + username)
	}

	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Username is required to register '# register <username>'\n")
	}

	username := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("User with the username %s already exists: %v\n", username, err)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("Something went wrong fetching existing users: %v\n", err)
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Something went wrong creating user %s: %v\n", username, err)
	}
	fmt.Println("> Created user")
	printUser(user)

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Something went wrong setting Config user %s: %v\n", username, err)
	}

	return nil
}

func handlerReset(s *State, cmd Command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("Something went wrong resetting the user table: %v\n", err)
	}
	return nil
}

func handlerUsers(s *State, cmd Command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Something went wrong getting users: %v\n", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUsername {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil

}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
