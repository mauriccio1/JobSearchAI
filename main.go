package main

import (
	"fmt"
	"jobsearch/internal/config"
	"jobsearch/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}
}

func main() {
	cfg := config.Load()
	if cfg.Name == "" || cfg.Certs == "" || cfg.Contact == "" || cfg.Education == nil || cfg.Port == "" {
		log.Fatal("Missing env variable - make sure are .env file or shell environment variables are set")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("starting server on port %s...\n", cfg.Port)
	if err := server.Start(sigChan, cfg); err != nil {
		log.Fatal(err)
	}

	
	
	
}
