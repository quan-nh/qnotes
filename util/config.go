package util

import (
	"encoding/json"
	"os"
)

type config struct {
	Port string
	Repo string
}

var Conf config

func LoadConfig() error {
	file, err := os.Open("conf.json")
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	return decoder.Decode(&Conf)
}
