package main

import (
	"fmt"
	"gator/internal/config"
	"os"
)

type state struct {
	con *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	callback map[string]func(*state, command) error
}

func (c *commands) run (s *state, cmd command) error {
	return c.callback[cmd.name](s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.callback[name] = f
}

func handlerLogin (s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("unexpected number of arguments: %v, command expects 1 argument", len(cmd.args))
	}

	err := s.con.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("received an error when trying to update username: %v", err)
	}

	fmt.Println("User has been successfuly set")
	return nil
}

func main() {
	configs, err := config.Read()
	if err != nil {
		fmt.Printf("Unexpected error occured when reading: %v\n", err)
	}

	running_state := state {
		con: &configs,
	}

	cmds := commands{
		make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Too few arguments were given")
		os.Exit(1)
	}

	cmd_name := args[1]
	var cmd_args []string = make([]string, 0)

	if len(args) > 2 {
		cmd_args = args[2:]
	}

	cmd := command{
		name: cmd_name,
		args: cmd_args,
	}

	err = cmds.run(&running_state, cmd)
	if err != nil {
		fmt.Printf("Unexpected error occured:\n\t%v\n", err)
		os.Exit(1)
	}
}