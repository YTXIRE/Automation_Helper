package main

import (
	"backend/internal/app/server"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := server.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(config)
	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
