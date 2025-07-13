package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

/*
======================================================

		Definitions

======================================================
*/

type state struct {
	cfg *config.Config
	db *database.Queries
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

/*
======================================================

		CLI commands

======================================================
*/

func handlerLogin (s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login command requires a username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])

	if err != nil {
		return  fmt.Errorf("user not found")
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("received an error when trying to update username: %v", err)
	}

	fmt.Println("User has been successfuly set")
	return nil
}

func handlerRegister (s *state, cmd command) error {
	if len(cmd.args) != 1  {
		return fmt.Errorf("register command requires a username")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		Name: cmd.args[0],
	})

	if err != nil {
		return fmt.Errorf("an unexpected error occured when creating user: %v", err)
	}

	s.cfg.SetUser(cmd.args[0])
	fmt.Println("User created successfuly:")
	fmt.Printf("\tID: %v | created_at: %v | updated_at: %v | name: %v\n", user.ID, user.CreatedAt, user.UpdatedAt.Time, user.Name)
	return nil
}

func handlerReset (s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %v", err)
	}

	return nil
}

/*
======================================================

		Entry point

======================================================
*/

func main() {
	var running_state state
	var configs config.Config
	var err error

	//		db connection
	{
		configs, err = config.Read()
		if err != nil {
			fmt.Printf("Unexpected error occured when reading config file: %v\n", err)
		}

		running_state = state {
			cfg: &configs,
		}

		db, err := sql.Open("postgres", running_state.cfg.Conn_str)
		if err != nil {
			fmt.Printf("Unexpected error occured when connecting to db: %v\n", err)
		}

		dbQueries := database.New(db)
		running_state.db = dbQueries
	}
	
	//		input handling
	{
		cmds := commands{
			make(map[string]func(*state, command) error),
		}

		cmds.register("login", handlerLogin)
		cmds.register("register", handlerRegister)
		cmds.register("reset", handlerReset)

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
}