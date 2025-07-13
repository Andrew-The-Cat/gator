package main

import (
	"fmt"
	"gator/internal/config"
)

func main() {
	configs, err := config.Read()
	if err != nil {
		fmt.Println("Unexpected error occured " + err.Error())
	}

	err = configs.SetUser("andrew")
	if err != nil {
		fmt.Println("Unexpected error occured " + err.Error())
	}

	configs, err = config.Read()
	if err != nil {
		fmt.Println("Unexpected error occured " + err.Error())
	}

	fmt.Printf("db_url: %v\ndb_user_name: %v", configs.Conn_str, configs.User_name)
}