package handler

import (
	"fmt"
	"net/http"

	"github.com/satyam-svg/resume-parser/internal/service"
)

func UploadProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	url, err := service.UploadToCloudinary(file, fileHeader)
	if err != nil {
		http.Error(w, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"image_url": "%s"}`, url)
}
