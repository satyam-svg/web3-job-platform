package routes

import (
	"net/http"
	"strings"

	"github.com/satyam-svg/resume-parser/internal/controller"
	"github.com/satyam-svg/resume-parser/internal/handler"
	"github.com/satyam-svg/resume-parser/internal/middleware"
	"github.com/satyam-svg/resume-parser/internal/service"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()

	// Resume Parsing APIs
	mux.HandleFunc("/upload", handler.UploadResumeHandler)
	mux.HandleFunc("/upload/profile-image", method("POST", handler.UploadProfileImageHandler))
	mux.HandleFunc("/signup", method("POST", controller.Signup))
	mux.HandleFunc("/login", method("POST", controller.Login))
	mux.HandleFunc("/reset-password", method("POST", controller.ResetPassword))
	mux.HandleFunc("/user/", method("GET", controller.GetUserByID))

	// Job APIs
	jobService := &service.JobService{DB: db}
	jobController := &controller.JobController{Service: jobService}

	// /jobs - POST: Create job | GET: List all jobs
	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			jobController.PostJob(w, r)
		case http.MethodGet:
			jobController.GetJobs(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// /jobs/... route handling
	mux.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/jobs/")

		switch {
		// GET /jobs/{jobID}/suggestions
		case strings.HasSuffix(path, "/suggestions") && r.Method == http.MethodGet:
			jobController.GetAISuggestions(w, r)
			return

		// GET /jobs/id/{jobID}
		case strings.HasPrefix(path, "id/") && r.Method == http.MethodGet:
			jobID := strings.TrimPrefix(path, "id/")
			jobController.GetJobByID(w, r, jobID)
			return

		// GET /jobs/recruiter/{recruiterID}
		case strings.HasPrefix(path, "recruiter/") && r.Method == http.MethodGet:
			recruiterID := strings.TrimPrefix(path, "recruiter/")
			jobController.GetJobsByRecruiterID(w, r, recruiterID)
			return

		default:
			http.NotFound(w, r)
		}
	})

	// AI Suggestions for a particular user
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/suggestions") && r.Method == http.MethodGet {
			userID := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/suggestions"), "/users/")
			jobController.GetUserAISuggestions(w, r, userID)
			return
		}
	})

	mux.HandleFunc("/credit/", jobController.GetUserCredit) // ✅ CORRECT

	mux.HandleFunc("/api/verify-payment", method("POST", controller.VerifyPaymentHandler))

	return middleware.CORS(mux)
}

// method ensures only a specific HTTP method is allowed
func method(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}
