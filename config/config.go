package config

import (
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
)

type Config struct {
	DB struct {
		Database string
		URL      string
	}

	Kafka struct {
		URL      string
		ClientID string
	}
}

func NewConfig(path string) *Config {
	c := new(Config)

	if f, err := os.Open(path); err != nil {
		log.Println("에러", "err", err.Error())
		return nil
	} else if err = toml.NewDecoder(f).Decode(c); err != nil {
		log.Println("에러", "err", err.Error())
		return nil
	} else {
		return c
	}
}
