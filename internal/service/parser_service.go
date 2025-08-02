package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/satyam-svg/resume-parser/config"
	"github.com/satyam-svg/resume-parser/internal/utils"
)

func ParseResume(filePath string) (string, error) {
	// 1. Extract text from PDF
	text, err := utils.ExtractTextFromPDF(filePath)
	if err != nil {
		return "", err
	}

	// 2. Sanitize raw PDF text
	sanitized := utils.SanitizeText(text)

	// 3. Prompt sent to Gemini
	prompt := fmt.Sprintf(`You are a resume parsing assistant. Extract the following information from the given resume and return it as valid JSON (DO NOT include "json" prefix):
{
  "full_name": "",
  "title": "",
  "experience": [
    {
      "company": "",
      "location": "",
      "title": "",
      "years": "",
      "description": ""
    }
  ],
  "education": [
    {
      "institution": "",
      "location": "",
      "degree": "",
      "gpa": "",
      "years": ""
    }
  ],
  "skills": [""],
  "location": "",
  "email": "",
  "phone": "",
  "current_company": "",
  "linkedin": "",
  "github": "",
  "portfolio": ""
}

Resume:
"%s"`, sanitized)

	// 4. Create API request
	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 5. Gemini API Key & Endpoint
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=%s", config.AppConfig.GeminiAPIKey)

	// 6. Make HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// 7. Handle non-200 errors
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Gemini API error: " + string(body))
	}

	// 8. Parse response JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}

	candidates, ok := parsed["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", errors.New("no candidates returned")
	}

	content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid content format")
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", errors.New("invalid parts format")
	}

	textOut, ok := parts[0].(map[string]interface{})["text"].(string)
	if !ok {
		return "", errors.New("no text found in response")
	}

	// 9. Clean markdown ```json ... ```
	cleaned := utils.CleanGeminiResponse(textOut)

	return cleaned, nil
}
