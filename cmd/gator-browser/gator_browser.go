package main

import _ "github.com/lib/pq"

import (
	"bootdevBlogAggerator/internal/cli"
	"fmt"
	"os"
)

func main() {
	state := cli.NewState()
	commands := cli.NewExplorerClient()

	cmd, err := cli.NewCommand(os.Args)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	err = commands.Run(state, cmd)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	os.Exit(0)

}
