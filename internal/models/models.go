package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jobsearch/internal/config"
	"jobsearch/internal/parser"
	"net/http"
	"strings"
)

const (
	ollamaUrl = "http://localhost:11434"
)

type OllamaRequest struct {
	Model   string          `json:"model"`
	Prompt  string          `json:"prompt"`
	System  string          `json:"system"`
	Stream  bool            `json:"stream"`
	Format  json.RawMessage `json:"format,omitempty"`
	Options Options         `json:"options"`
}

type Options struct {
	Temperature float32 `json:"temperature"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type ExperienceResult struct {
	Experience []parser.Job `json:"experience"`
	err        error
}

type SummaryResult struct {
	Summary string `json:"summary"`
	err     error
}

func CheckOllama() error {
	resp, err := http.Get(ollamaUrl + "/api/version")
	if err != nil {
		return fmt.Errorf("health check failed to run: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama server is not available - please make sure ollama is running")
	}
	return nil
}

func GenerateEnv(resume string) (string, error) {
	req := &OllamaRequest{
		Model:   "llama3.2:3b",
		Prompt:  fmt.Sprintf(setupPrompt, resume),
		System:  setupSysPrompt,
		Stream:  false,
		Options: Options{Temperature: 0.1},
	}
	response, err := generate(req)
	if err != nil {
		return "", err
	}
	return response, nil
}

func TrimJD(fullJD string, req *OllamaRequest) (string, error) {
	if req == nil {
		req = &OllamaRequest{
			Model:   "llama3.2:3b",
			Prompt:  jdPrompt + fullJD,
			System:  jdSystemPrompt,
			Stream:  false,
			Options: Options{Temperature: 0.1},
		}
	}
	response, err := generate(req)
	if err != nil {
		return "", err
	}
	return response, nil
}

// RewriteAndStructure replaces the old RewriteResume + StructureResume flow.
// It runs summary and experience rewrites concurrently, then structures the result.
func RewriteAndStructure(resume, trimmedJD string, cfg *config.Config) (*parser.Resume, error) {
	summCh := make(chan *SummaryResult, 1)
	expCh := make(chan *ExperienceResult, 1)

	// rewrite summary (one call)
	go func() {
		summary, err := rewriteSummary(resume, trimmedJD)
		if err != nil {
			summCh <- &SummaryResult{err: err}
			return
		}
		summCh <- &SummaryResult{Summary: summary}
	}()

	// rewrite experience then structure it (two sequential calls)
	go func() {
		rewritten, err := rewriteExperience(resume, trimmedJD)
		if err != nil {
			expCh <- &ExperienceResult{err: err}
			return
		}
		exp, err := extractExp(rewritten)
		if err != nil {
			expCh <- &ExperienceResult{err: err}
			return
		}
		expCh <- exp
	}()

	summary := <-summCh
	exp := <-expCh

	if summary.err != nil {
		return nil, fmt.Errorf("RewriteSummary failed: %w", summary.err)
	}
	if exp.err != nil {
		return nil, fmt.Errorf("RewriteExperience failed: %w", exp.err)
	}

	return &parser.Resume{
		Name:       cfg.Name,
		Contact:    cfg.Contact,
		Summary:    summary.Summary,
		Experience: exp.Experience,
		Skills:     cfg.Skills,
		Certs:      cfg.Certs,
		Education:  cfg.Education,
	}, nil
}

func rewriteSummary(resume, trimmedJD string) (string, error) {
	req := &OllamaRequest{
		Model:   "gemma2:9b",
		Prompt:  fmt.Sprintf(rewriteSummaryPrompt, resume, trimmedJD),
		System:  rewriteSummarySysPrompt,
		Stream:  false,
		Options: Options{Temperature: 0.2},
	}
	response, err := generate(req)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response), nil
}

func rewriteExperience(resume, trimmedJD string) (string, error) {
	req := &OllamaRequest{
		Model:   "gemma2:9b",
		Prompt:  fmt.Sprintf(rewriteExperiencePrompt, resume, trimmedJD),
		System:  rewriteExperienceSysPrompt,
		Stream:  false,
		Options: Options{Temperature: 0.2},
	}
	response, err := generate(req)
	if err != nil {
		return "", err
	}
	return response, nil
}

func extractExp(text string) (*ExperienceResult, error) {
	req := &OllamaRequest{
		Model:   "gemma2:9b",
		Prompt:  text,
		System:  extractExperienceSysPrompt,
		Stream:  false,
		Format:  json.RawMessage(experienceSchema),
		Options: Options{Temperature: 0.1},
	}
	response, err := generate(req)
	if err != nil {
		return nil, err
	}
	exp := &ExperienceResult{}
	if err := json.Unmarshal([]byte(response), exp); err != nil {
		return nil, fmt.Errorf("extractExp failed to unmarshal: %w", err)
	}
	return exp, nil
}

func DryRun() *parser.Resume {
	return &parser.Resume{
		Name:    "John Doe",
		Contact: "Miami, FL | 555-123-4567 | johndoe@gmail.com | linkedin.com/in/johndoe",
		Summary: "Platform engineer with 4 years of experience building Go services and cloud infrastructure on AWS and GCP.",
		Certs:   "Cisco CCNA | CompTIA Security+ | CompTIA Network+",
		Education: []string{
			"B.S. Computer Science, Florida International University",
			"College Certificate, Network Security, Miami Dade College",
		},
		Skills: []parser.Skill{
			{Label: "Languages", Value: "Go, Python, Bash, PowerShell"},
			{Label: "Cloud & Platform", Value: "AWS, GCP, Docker, Kubernetes"},
			{Label: "IaC & CI/CD", Value: "Terraform, Pulumi, GitHub Actions"},
			{Label: "Networking", Value: "BGP, VLANs, firewall design, VPN"},
			{Label: "Observability", Value: "Prometheus, Grafana, Loki"},
			{Label: "Identity & Security", Value: "Okta, Active Directory, RBAC"},
		},
		Experience: []parser.Job{
			{
				Title:   "Platform Engineer",
				Company: "Acme Corp",
				Dates:   "Jan 2023 – Present",
				Sections: []parser.JobSection{
					{
						Header: "Infrastructure & CI/CD",
						Bullets: []string{
							"Built and maintained CI/CD pipelines using GitHub Actions across 10+ services.",
							"Standardized Terraform IaC conventions adopted across all production services.",
						},
					},
				},
			},
		},
	}
}

func generate(reqBody *OllamaRequest) (string, error) {
	url := "http://localhost:11434/api/generate"
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("generate failed to marshal request: %w", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("generate failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("generate failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("generate failed with status: %s", resp.Status)
	}
	response := &OllamaResponse{}
	if err := json.Unmarshal(respData, response); err != nil {
		return "", fmt.Errorf("generate failed to parse response: %w", err)
	}
	if !response.Done {
		return "", fmt.Errorf("generate failed: generation didn't finish")
	}
	return response.Response, nil
}
