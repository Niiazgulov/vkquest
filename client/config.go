package client

import (
	"flag"
	"os"
)

const (
	baseURLAddr = "http://localhost:8080/"
)

type Config struct {
	BaseURLAddress string `json:"base_url"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		BaseURLAddress: "",
	}

	flag.StringVar(&cfg.BaseURLAddress, "b", "", "base url")
	flag.Parse()
	cfg.BaseURLAddress = pickFirstNonEmpty(cfg.BaseURLAddress, os.Getenv("BASE_URL"), baseURLAddr)

	return cfg, nil
}

var Cfg Config

func pickFirstNonEmpty(strings ...string) string {
	for _, str := range strings {
		if str != "" {
			return str
		}
	}
	return ""
}
