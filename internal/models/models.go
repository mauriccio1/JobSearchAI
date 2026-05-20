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
