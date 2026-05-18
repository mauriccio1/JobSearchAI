package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	jdSystemPrompt = "You extract job requirements from job descriptions. You remove everything that is not relevant to evaluating a candidate's qualifications. You output only the extracted sections with no commentary."
	jdPrompt = `Extract only the following from this job description. Output nothing else.
JOB TITLE
KEY RESPONSIBILITIES (bullet points only, no prose)
REQUIRED SKILLS AND QUALIFICATIONS (bullet points only)
PREFERRED/DESIRED SKILLS (bullet points only, if listed)
TECHNOLOGIES MENTIONED (flat comma-separated list)
Remove everything else: company description, benefits, compensation, legal disclaimers, interview process, diversity statements, travel requirements, office location details.
---JOB DESCRIPTION---
	`
	resumeRewriteSysPrompt = "You are a resume rewriting engine. You receive a resume and a job description. You rewrite the resume to better match the job description. You ONLY reword and reorder what already exists. You NEVER add or remove information. Every skill, certification, employer, bullet point, and education entry from the original must appear in your output. Your output is always a complete resume with no commentary."

	resumeRewritePrompt = `Rewrite my resume to match the job description. Follow these rules exactly.
RULES:

OUTPUT: Only the final resume. No commentary, no analysis, no tips, no questions.
STRUCTURE — use this exact section order:
SUMMARY
EXPERIENCE
TECHNICAL SKILLS
CERTIFICATIONS
EDUCATION
EXPERIENCE RULES:
Keep ALL four employers in this exact order:
a) Systems & Cloud Engineer | The Pharmacy Hub | Feb 2025 – May 2026
b) Support Engineer | Ringlogix (White Label VoIP) | Jun 2024 – Jan 2025
c) Server & Systems Support Specialist | World Wide Tech Services | Jun 2023 – May 2024
d) IT Support Technician | Cook Technology Corp. | Dec 2022 – May 2023
Do NOT merge, rename, or remove any employer.
Do NOT move bullet points between employers.
You may reword bullets to use keywords from the job description.
You may reorder bullets within each employer to put the most relevant first.
SUMMARY: Rewrite using keywords from the job description. Use ONLY skills and experience from my original resume.
TECHNICAL SKILLS — include ALL of the following, reordered for relevance to the JD:
Languages: Go (primary), Python, Bash, PowerShell. Working knowledge of TypeScript/JavaScript.
Cloud & Platform: GCP (Cloud Run, Vertex AI, Managed AD, IAM, Cloud VPN), AWS, Docker, Kubernetes, Proxmox.
IaC & CI/CD: Pulumi (Go), Terraform, Bitbucket Pipelines, GitHub Actions, Workload Identity Federation.
Networking: BGP, HA VPN, VLANs, firewall design, network segmentation, zero-trust (WireGuard).
Observability: GCP Cloud Logging, Grafana, Prometheus, Loki.
Identity & Security: Okta (SSO, SAML, SCIM), Active Directory, Google Workspace, RBAC, HIPAA controls.
CERTIFICATIONS — print exactly:
Cisco CCNA, CompTIA CySA+, CompTIA Security+, CompTIA Network+, CompTIA A+
EDUCATION — print exactly:
A.S. Cybersecurity, Miami Dade College
College Certificate, Network Security, Miami Dade College
DO NOT ADD anything not in the original resume. No new skills, tools, certifications, degrees, or employers.
DO NOT REMOVE anything from the original resume. Every item must appear.
---RESUME---
`
)

type OllamaRequest struct {
    Model   string  `json:"model"`
    Prompt  string  `json:"prompt"`
    System  string  `json:"system"`
    Stream  bool    `json:"stream"`
    Options Options `json:"options"`
}
type Options struct {
	Temperature float32 `json:"temperature"` 
}

type OllamaResponse struct {
    Response string `json:"response"`
    Done     bool   `json:"done"`
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

