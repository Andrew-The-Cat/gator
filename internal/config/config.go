package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Conn_str 	string 	`json:"db_url"`
	User_name 	string 	`json:"current_user_name"`
}

func get_gator_path() string {
	home_path, _ := os.UserHomeDir()
	return home_path + "/.gatorconfig.json"
}

func Read() (Config, error) {
	json_data, err := os.ReadFile( get_gator_path() )

	if err != nil {
		return Config{}, err
	}

	var returned Config

	if err := json.Unmarshal(json_data, &returned); err != nil {
		return Config{}, err
	}
	return returned, nil
}

func (c Config) SetUser(user_name string) error {
	c.User_name = user_name

	json_data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(get_gator_path(), json_data, os.ModeDevice)
	if err != nil {
		return err
	}

	return nil
}