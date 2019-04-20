package main

import (
	"cookbook/config"
	"cookbook/server"
	"cookbook/utils"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	environment := flag.String("e",
		utils.GetEnv("ENVIRONMENT", "dev"), "Sets the environment and pulls in config accordingly.")

	flag.Usage = func() {
		log.Fatal("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	// Configuration gets initialized here
	err := config.Init(*environment)

	if err != nil {
		log.Error("Failed to set environment configs via file.")
	}

	// Initialize server
	server.Init()

}
