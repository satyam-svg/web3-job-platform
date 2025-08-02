package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// MatchResult represents a structured match response
type MatchResult struct {
	Matches []struct {
		Name          string `json:"name"`
		Email         string `json:"email"`
		MatchingScore int    `json:"matching_score"`
		Reasoning     string `json:"reasoning"`
		Recommended   bool   `json:"recommended"`
	} `json:"matches"`
}

type RecommendationResult struct {
	Recommendations []struct {
		Title         string `json:"title"`
		Company       string `json:"company"`
		MatchingScore int    `json:"matching_score"`
		Reasoning     string `json:"reasoning"`
		Recommended   bool   `json:"recommended"`
	} `json:"recommendations"`
}

// CallGemini sends a prompt to Gemini and returns structured candidate suggestions
func CallGemini(prompt string) (*MatchResult, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("❌ GEMINI_API_KEY not set in environment")
	}

	// Payload structure
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to marshal payload: %v", err)
	}

	// Make request to Gemini
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=%s", apiKey)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("❌ Request failed: %v", err)
	}
	defer res.Body.Close()

	// Parse Gemini API response
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("❌ Failed to decode response: %v", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("❌ Gemini returned no content")
	}

	rawText := result.Candidates[0].Content.Parts[0].Text
	fmt.Println("📥 Gemini Raw Response:\n", rawText)

	// Try to extract JSON block from markdown/code
	jsonStr := extractJSONFromText(rawText)

	var parsed MatchResult
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("❌ Failed to parse structured JSON: %v\nRaw: %s", err, jsonStr)
	}

	return &parsed, nil
}

// extractJSONFromText extracts the first JSON-like code block from a Gemini text response
func extractJSONFromText(text string) string {
	// Remove markdown code fences if present
	re := regexp.MustCompile("(?s)```json(.*?)```")
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	// Fallback if no ```json block found
	return strings.TrimSpace(text)
}

func CallGeminiForRecommendations(prompt string) (*RecommendationResult, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("❌ GEMINI_API_KEY not set in environment")
	}

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to marshal payload: %v", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=%s", apiKey)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("❌ Request failed: %v", err)
	}
	defer res.Body.Close()

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("❌ Failed to decode response: %v", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("❌ Gemini returned no content")
	}

	rawText := result.Candidates[0].Content.Parts[0].Text
	fmt.Println("📥 Gemini Raw Response:\n", rawText)

	jsonStr := extractJSONFromText(rawText)

	var parsed RecommendationResult
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("❌ Failed to parse recommendations JSON: %v\nRaw: %s", err, jsonStr)
	}

	return &parsed, nil
}
