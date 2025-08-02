package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/satyam-svg/resume-parser/internal/service"
)

func UploadResumeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "resume-*.pdf")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	tempFile.Close()

	// Call the Gemini service
	jsonOutput, err := service.ParseResume(tempFile.Name())
	if err != nil {
		http.Error(w, "Failed to parse resume: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, jsonOutput)
}
