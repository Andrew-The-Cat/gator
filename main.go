package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"
	"time"
	"gator/internal/rss"

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
		return fmt.Errorf("command requires a username")
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
	fmt.Println("Successfuly created user:")
	fmt.Printf("\tID: %v | created_at: %v | updated_at: %v | name: %v\n", user.ID, user.CreatedAt, user.UpdatedAt.Time, user.Name)
	return nil
}

func handlerReset (s *state, cmd command) error {
	err := s.db.UsersReset(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %v", err)
	}

	err = s.db.FeedsReset(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete feeds: %v", err)
	}

	err = s.db.FeedFollowsReset(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete feed followdL: %v", err)
	}

	return nil
}

func handlerAgg (s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %v", err)
	}

	fmt.Println(feed)
	return nil
}

func handlerUsers (s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("* %v", user.Name)
		if user.Name == s.cfg.User_name {
			fmt.Printf(" (current)")
		}

		fmt.Print("\n")
	}

	return nil
}

func handlerAddFeed (s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("command requires a name and a url")
	}

	res, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding the feed to the database: %v", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: res.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %v", err)
	}

	fmt.Println("Successfuly added feed:")
	fmt.Printf("\tname: %v | url: %v\n", res.Name, res.Url)
	return nil
}

func handlerFeeds (s *state, cmd command) error {
	data, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving feeds: %v", err)
	}

	for _, row := range data {
		fmt.Printf(" * %v - %v: %v\n", row.UserName, row.Name, row.Url)
	}

	return nil
}

func handlerFollow (s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("command requires the url of the feed you want to follow")
	}

	targetFeed, err := s.db.GetFeedByUrl(context.Background() ,cmd.args[0])
	if err != nil {
		return fmt.Errorf("error retrieving requested feed: %v", err)
	}

	response, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: targetFeed.ID,
	})

	if err != nil {
		return fmt.Errorf("error creating feed follow: %v", err)
	}

	fmt.Printf("Successfuly followed feed for user %v:\n", s.cfg.User_name)
	fmt.Printf("\t* name: %v | url: %v\n", response.FeedName, cmd.args)
	return nil
}

func handlerFollowing (s *state, cmd command, user database.User) error {
	data, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error retrieving feed follows: %v", err)
	}

	for _, row := range data {
		fmt.Printf(" * %v: %v\n", row.FeedName, row.FeedUrl)
	}
	return nil
}

func handlerUnfollow (s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("command requires the url of the feed you want to unfollow")
	}

	err := s.db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserID: user.ID,
		Url: cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("error deleting feed follow: %v", err)
	}

	fmt.Println("Successfully unfollowed feed")
	return nil
}

/*
======================================================

		Function wrappers

======================================================
*/

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func (s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.User_name)
		if err != nil {
			return fmt.Errorf("error retrieving current user: %v", err)
		}

		return handler(s, cmd, user)
	}
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
		cmds.register("users", handlerUsers)
		cmds.register("agg", handlerAgg)
		cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
		cmds.register("feeds", handlerFeeds)
		cmds.register("follow", middlewareLoggedIn(handlerFollow))
		cmds.register("following", middlewareLoggedIn(handlerFollowing))
		cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

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