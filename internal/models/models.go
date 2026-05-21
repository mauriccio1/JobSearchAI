package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jobsearch/internal/config"
	"jobsearch/internal/parser"
	"net/http"
)


type OllamaRequest struct {
    Model   string  `json:"model"`
    Prompt  string  `json:"prompt"`
    System  string  `json:"system"`
    Stream  bool    `json:"stream"`
	Format  json.RawMessage `json:"format,omitempty"`
    Options Options `json:"options"`
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
	err 		error
}

type MetaResult struct {
    Name      string   `json:"name"`
    Contact   string   `json:"contact"`
    Summary   string   `json:"summary"`
    Skills    []parser.Skill  `json:"skills"`
    Certs     string   `json:"certs"`
    Education []string `json:"education"`
	err       error
}

func TrimJD(fullJD string, req *OllamaRequest) (string, error) {
	if req == nil {
		req = &OllamaRequest{
			Model: "llama3.2:3b",
			Prompt: jdPrompt + fullJD,
			System: jdSystemPrompt,
			Stream: false,
			Options: Options{
				Temperature: 0.1,
			},
		}
	}

	response, err := generate(req)
	if err != nil {
		return "", err
	}

	return response, nil
}

func RewriteResume(resume, trimmedJD string, req *OllamaRequest) (string, error) {
	if req == nil {
		req = &OllamaRequest{
			Model: "gemma2:9b",
			Prompt: resumeRewritePrompt + resume + "\n---JOB DESCRIPTION---\n" + trimmedJD,
			System: resumeRewriteSysPrompt,
			Stream: false,
			Options: Options{
				Temperature: 0.2,
			},
		}
	}
	response, err := generate(req)
	if err != nil {
		return "", err
	}

	return response, nil

}

func StructureResume(text string, cfg *config.Config) (*parser.Resume, error) {
	expCh := make(chan *ExperienceResult, 1)
    metaCh := make(chan *MetaResult, 1)

	newResume := &parser.Resume{}

	
	go func() {
		exp, err := ExtractExp(text)
		if err != nil {
			expCh <-&ExperienceResult{err: err}
			return 
		}
		expCh <-exp

	}()

	go func() {
		meta, err := ExtractMetaData(text)
		if err != nil {
			metaCh <-&MetaResult{err: err}
			return 
		}
		
		metaCh <-meta
	}()
	
	exp := <-expCh
	meta := <-metaCh

	if exp.err != nil {
    	return nil, fmt.Errorf("ExtractExp failed: %w", exp.err)
    }
    if meta.err != nil {
    	return nil, fmt.Errorf("ExtractMetaData failed: %w", meta.err)
    }


	

	newResume.Experience = exp.Experience
	newResume.Name = cfg.Name
	newResume.Contact = cfg.Contact
	newResume.Summary = meta.Summary
	newResume.Skills =  meta.Skills
	newResume.Certs =  cfg.Certs
	newResume.Education = cfg.Education

	return newResume, nil

} 

func ExtractExp(text string) (*ExperienceResult, error) {
	req := &OllamaRequest{
		Model: "gemma2:9b",
		Prompt: text,
		System: extractExperienceSysPrompt,
		Stream: false,
		Format: json.RawMessage(experienceSchema),
		Options: Options{
			Temperature: 0.1,
		},
	}

	response, err := generate(req)
	if err != nil {
		return nil, err
	}

	exp := &ExperienceResult{}
	if err := json.Unmarshal([]byte(response), exp); err != nil {
		return nil, fmt.Errorf("ExperienceResult failed to unmarshal response: %w", err)
	}
	return exp, nil

}

func ExtractMetaData(text string) (*MetaResult, error) {
	req := &OllamaRequest{
		Model: "gemma2:9b",
		Prompt: text,
		System: extractMetaSysPrompt,
		Stream: false,
		Format: json.RawMessage(metaSchema),
		Options: Options{
			Temperature: 0.1,
		},
	}

	response, err := generate(req)
	if err != nil {
		return nil, err
	}

	meta := &MetaResult{}
	if err := json.Unmarshal([]byte(response), meta); err != nil {
		return nil, fmt.Errorf("ExtractMetaData failed to unmarshal response: %w", err)
	}
	return meta, nil

}


func DryRun() *parser.Resume {
    return &parser.Resume{
        Name:    "John Doe",
        Contact: "Miami, FL | 555-123-4567 | johndoe@gmail.com | linkedin.com/in/johndoe",
        Summary: "Platform engineer with 4 years of experience building Go services and cloud infrastructure on AWS and GCP. Shipped internal developer tooling, CI/CD pipelines, and IaC with Terraform and Pulumi. CCNA certified with strong networking foundations.",
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
                Intro:   "Platform engineer on a four-person infrastructure team supporting a 300-employee operation.",
                Sections: []parser.JobSection{
                    {
                        Header: "Infrastructure & CI/CD",
                        Bullets: []string{
                            "Built and maintained CI/CD pipelines using GitHub Actions across 10+ services.",
                            "Standardized Terraform IaC conventions adopted across all production services.",
                            "Implemented Workload Identity Federation for keyless GCP authentication.",
                        },
                    },
                    {
                        Header: "Platform Services",
                        Bullets: []string{
                            "Shipped internal Go HTTP service on Cloud Run as abstraction layer over Jira REST API.",
                            "Built employee lifecycle automation reducing onboarding from 30 minutes to under 1 minute.",
                        },
                    },
                },
            },
            {
                Title:   "DevOps Engineer",
                Company: "Beta Systems",
                Dates:   "Jun 2021 – Dec 2022",
                Intro:   "DevOps engineer supporting cloud migration and infrastructure automation.",
                Sections: []parser.JobSection{
                    {
                        Header: "Cloud Migration",
                        Bullets: []string{
                            "Migrated 15 legacy services from on-premise to AWS EC2 and ECS.",
                            "Reduced infrastructure costs by 30% through right-sizing and reserved instances.",
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
		return "", fmt.Errorf("generate failed to parse response: %w", err)
	}
	resp, err := http.Post(url, "application/json",bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("generate failed: %w", err)
	}
	defer func() {_ = resp.Body.Close()}()
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

func GeneratePDF() {


}
