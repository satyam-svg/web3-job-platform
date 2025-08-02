// controller/job_controller.go
package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/satyam-svg/resume-parser/internal/model"
	"github.com/satyam-svg/resume-parser/internal/service"
	"github.com/satyam-svg/resume-parser/internal/utils"
)

type JobController struct {
	Service *service.JobService
}

// Create a new job
func (jc *JobController) PostJob(w http.ResponseWriter, r *http.Request) {
	var job model.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	job.CreatedAt = time.Now()

	if err := jc.Service.CreateJob(&job); err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(job)
}

// Get all jobs
func (jc *JobController) GetJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := jc.Service.GetJobs()
	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(jobs)
}

// Get jobs by recruiter ID
func (jc *JobController) GetJobsByRecruiterID(w http.ResponseWriter, r *http.Request, recruiterID string) {
	jobs, err := jc.Service.GetJobsByRecruiter(recruiterID)
	if err != nil {
		http.Error(w, "Failed to fetch jobs for recruiter", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(jobs)
}

// Get AI Suggestions using Gemini
func (jc *JobController) GetAISuggestions(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/jobs/")
	jobIDStr := strings.TrimSuffix(path, "/suggestions")

	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	job, err := jc.Service.GetJobByID(uint(jobID))
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	recruiter, err := jc.Service.GetUserByID(job.RecruiterID.String())
	if err != nil {
		http.Error(w, "Recruiter not found", http.StatusNotFound)
		return
	}

	if err := deductCredit(recruiter, jc.Service); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	users, err := jc.Service.GetRelevantStudents()
	if err != nil {
		http.Error(w, "Failed to fetch students", http.StatusInternalServerError)
		return
	}

	prompt := buildPrompt(*job, users)
	matchResult, err := utils.CallGemini(prompt)
	if err != nil {
		http.Error(w, "Gemini failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(matchResult)
}

// Prompt Builder
func buildPrompt(job model.Job, users []model.User) string {
	prompt := fmt.Sprintf(`You are an AI recruitment assistant.
Analyze the job below and match it with the best candidates.

Job Profile:
Title: %s
Company: %s
Location: %s
Tags: %s
Description: %s

Candidates:
`, job.Title, job.Company, job.Location, job.Tags, job.Description)

	for _, u := range users {
		prompt += fmt.Sprintf(`
Name: %s
Email: %s
Location: %s
Skills: %s
Experience: %s
`, u.FullName, u.Email, u.Location, u.Skills, u.Experience)
	}

	prompt += `
Respond only with a valid JSON in the following format:
{
  "matches": [
    {
      "name": "Candidate Name",
      "email": "email@example.com",
      "matching_score": 0-100,
      "reasoning": "explanation",
      "recommended": true/false
    },
    ...
  ]
}`

	return prompt
}

// Get a single job by ID
func (jc *JobController) GetJobByID(w http.ResponseWriter, r *http.Request, id string) {
	jobID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	job, err := jc.Service.GetJobByID(uint(jobID))
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(job)
}

// Get AI Job Suggestions for a Specific User
func (jc *JobController) GetUserAISuggestions(w http.ResponseWriter, r *http.Request, userID string) {
	user, err := jc.Service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := deductCredit(user, jc.Service); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	jobs, err := jc.Service.GetAllJobs()
	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}

	// Build prompt
	prompt := fmt.Sprintf(`You are an AI career advisor.
Suggest jobs for the following user based on their profile.

User Profile:
Name: %s
Email: %s
Location: %s
Skills: %s
Experience: %s

Jobs:
`, user.FullName, user.Email, user.Location, user.Skills, user.Experience)

	for _, job := range jobs {
		prompt += fmt.Sprintf(`
Title: %s
Company: %s
Location: %s
Tags: %s
Description: %s
`, job.Title, job.Company, job.Location, job.Tags, job.Description)
	}

	prompt += `
Respond only with a valid JSON like this:
{
  "recommendations": [
    {
      "title": "Job Title",
      "company": "Company Name",
      "matching_score": 0-100,
      "reasoning": "Why it's a match",
      "recommended": true/false
    }
  ]
}`

	matchResult, err := utils.CallGeminiForRecommendations(prompt)
	if err != nil {
		http.Error(w, "Gemini failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(matchResult.Recommendations)
}

func deductCredit(user *model.User, service *service.JobService) error {
	if user.Credits <= 0 {
		return fmt.Errorf("You don't have enough credits")
	}
	user.Credits -= 1
	return service.UpdateUserCredits(user.ID, user.Credits)
}

// GetUserCredit returns the remaining credits of a user
func (jc *JobController) GetUserCredit(w http.ResponseWriter, r *http.Request) {
	vars := strings.Split(r.URL.Path, "/")
	if len(vars) < 3 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userID := vars[2]

	user, err := jc.Service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": user.ID,
		"credits": user.Credits,
	})
}
