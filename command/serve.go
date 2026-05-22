package command

import (
    "fmt"
    "jobsearch/internal/config"
    "jobsearch/internal/server"
    "log"
    "os"
    "os/signal"
    "syscall"
    "github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the resume generation server",
    Run: func(cmd *cobra.Command, args []string) {
        cfg := config.Load()
        if cfg.Name == "" || cfg.Certs == "" || cfg.Contact == "" || cfg.Education == nil || cfg.Port == "" {
            log.Fatal("missing required config — run 'jobsearchai setup' first")
        }
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        fmt.Printf("starting server on port %s...\n", cfg.Port)
        if err := server.Start(sigChan, cfg); err != nil {
            log.Fatal(err)
        }
    },
}
