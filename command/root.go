package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "JobSearchAI",
	Short: "AI-powered resume optimization tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
	rootCmd.AddCommand(serveCmd)
    rootCmd.AddCommand(setupCmd)
}
