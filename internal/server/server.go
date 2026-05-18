package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jobsearch/internal/models"
	"log"
	"net/http"
	"os"
	"time"
)
type RewriteRequest struct {
    JD string `json:"jd"`
}

func Start(sigChan <- chan os.Signal) error {
	resume, err := os.ReadFile("./resume/base_resume.txt")
		if err != nil {
		return fmt.Errorf("failed to read resume text file: %w", err)
		}

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

		newResume, err := models.RewriteResume(string(resume), newJD, nil)	
		if err != nil {
		log.Printf("RewriteResume failed: %v", err)
        http.Error(w, "failed to rewrite resume", http.StatusInternalServerError)
        return
    	}		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"resume": newResume})

	}))
	
	

	server := http.Server{
		Addr: ":8080",
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
