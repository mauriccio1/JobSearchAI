package main

import (
	"fmt"
	"jobsearch/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("starting server on port 8080...")
	if err := server.Start(sigChan); err != nil {
		log.Fatal(err)
	}

	
	
	
}
