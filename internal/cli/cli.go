package cli

import (
	"bootdevBlogAggerator/internal/config"
	"bootdevBlogAggerator/internal/database"
	"database/sql"
	"fmt"
	"strings"
)

type State struct {
	cfg *config.Config
	db  *database.Queries
}

func NewState() *State {
	newConfig, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", newConfig.DBUrl)
	dbQueries := database.New(db)

	state := State{cfg: &newConfig, db: dbQueries}
	return &state
}

type Command struct {
	Name string
	Args []string
}

func NewCommand(args []string) (Command, error) {
	if len(args) < 2 {
		return Command{}, fmt.Errorf("something went wrong")
	}

	name := args[1]
	args = args[2:]

	return Command{name, args}, nil
}

type Client struct {
	Commands map[string]func(*State, Command) error
}

func (c *Client) Register(name string, f func(*State, Command) error) {
	name = strings.ToLower(name)
	if _, exists := c.Commands[name]; exists {
		fmt.Printf("command '%s' was registered and will be overwritten", name)
	}
	c.Commands[name] = f
}

func (c *Client) Run(s *State, cmd Command) error {

	callback, exists := c.Commands[cmd.Name]
	if !exists {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}
	err := callback(s, cmd)
	if err != nil {
		return err
	}

	return nil
}
