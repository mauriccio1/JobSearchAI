package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jobsearch/internal/config"
	"jobsearch/internal/models"
	"jobsearch/internal/pdf"
	"log"
	"net/http"
	"os"
	"time"
)
type RewriteRequest struct {
    JD string `json:"jd"`
}

func Start(sigChan <- chan os.Signal, cfg *config.Config) error {
	resume, err := os.ReadFile("./resume/base_resume.txt")
		if err != nil {
		return fmt.Errorf("failed to read resume text file: %w", err)
		}

	http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
 		w.WriteHeader(http.StatusOK)
	}))


	http.Handle("/api/resume/dry-run", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pdfBytes, err := pdf.Generate(models.DryRun())
		if err != nil {
			log.Printf("GeneratePDF failed: %v", err)
    		http.Error(w, "failed to generate pdf", http.StatusInternalServerError)
    		return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")
		w.Write(pdfBytes)


	}))



	http.Handle("/api/resume/rewrite", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "invalid method", http.StatusBadRequest)
			return
		}
		defer func() {_ = r.Body.Close()}()
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read request body: %v", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	
		rewriteReq := &RewriteRequest{}
		if err := json.Unmarshal(reqBody, rewriteReq); err != nil {
			log.Printf("failed to unmarshal request body: %v", err)

			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		newJD, err := models.TrimJD(rewriteReq.JD, nil)		
		if err != nil {
			log.Printf("TrimJD failed: %v", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		newResume, err := models.RewriteAndStructure(string(resume), newJD, cfg)	
		if err != nil {
		log.Printf("RewriteResume failed: %v", err)
        http.Error(w, "failed to rewrite resume", http.StatusInternalServerError)
        return
    	}

		pdfBytes, err := pdf.Generate(newResume)	
		if err != nil {
			log.Printf("GeneratePDF failed: %v", err)
    		http.Error(w, "failed to generate pdf", http.StatusInternalServerError)
    		return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")
		w.Write(pdfBytes)
	}))
	
	

	server := http.Server{
		Addr: ":" + cfg.Port,
	}
	

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("server error; %s", err)
		}
	}()
	<-sigChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}
