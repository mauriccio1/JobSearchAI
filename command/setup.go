package command

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "jobsearch/internal/models"
    "github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
    Use:   "setup",
    Short: "Read your resume and generate .env",
    Run: func(cmd *cobra.Command, args []string) {
        // 1. check resume exists
        resumeBytes, err := os.ReadFile("resume/base_resume.txt")
        if err != nil {
            log.Fatal("resume/base_resume.txt not found — create it and re-run setup")
        }

        // 2. check ollama
        if err := models.CheckOllama(); err != nil {
            log.Fatalf("Ollama unreachable: %s", err)
        }

        fmt.Println("reading resume...")

        // 3. call LLM
        result, err := models.GenerateEnv(string(resumeBytes))
        if err != nil {
            log.Fatal("failed to generate .env: ", err)
        }

        // 4. show result
        fmt.Println("\n--- generated .env ---")
        fmt.Println(result)
        fmt.Println("----------------------")

        // 5. confirm
        fmt.Print("\nwrite this to .env? (y/n): ")
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        answer := strings.TrimSpace(strings.ToLower(scanner.Text()))

        if answer != "y" {
            fmt.Println("edit resume/base_resume.txt and re-run setup")
            return
        }

        // 6. write .env
        if err := os.WriteFile(".env", []byte(result), 0600); err != nil {
            log.Fatal("failed to write .env: ", err)
        }

        fmt.Println(".env written. run 'go run main.go serve' to start.")
    },
}
