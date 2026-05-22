// main.go
package main

import (
    "jobsearch/command"
    "github.com/joho/godotenv"
    "log"
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("no .env file found, using environment variables")
    }
}

func main() {
    command.Execute()
}
