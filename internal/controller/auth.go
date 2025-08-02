package controller

import (
	"encoding/json"
	"net/http"

	"github.com/satyam-svg/resume-parser/config"
	"github.com/satyam-svg/resume-parser/internal/model"
	"github.com/satyam-svg/resume-parser/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// ---------- Request Structs ----------
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email      string             `json:"email"`
	Password   string             `json:"password"`
	FullName   string             `json:"full_name"`
	Title      string             `json:"title"`
	Location   string             `json:"location"`
	Phone      string             `json:"phone"`
	Company    string             `json:"current_company"`
	LinkedIn   string             `json:"linkedin"`
	GitHub     string             `json:"github"`
	Portfolio  string             `json:"portfolio"`
	Skills     string             `json:"skills"`
	Image      string             `json:"image"`
	Role       string             `json:"role"` // recruiter, applicant, admin
	Education  []model.Education  `json:"education"`
	Experience []model.Experience `json:"experience"`
}

// ---------- Filtered Response ----------
func filterUserResponse(user model.User) interface{} {
	if user.Role == "recruiter" {
		var jobCount int64
		config.DB.Model(&model.Job{}).Where("recruiter_id = ?", user.ID).Count(&jobCount)

		return model.RecruiterResponse{
			ID:              user.ID,
			FullName:        user.FullName,
			Email:           user.Email,
			Phone:           user.Phone,
			Image:           user.Image,
			CurrentCompany:  user.CurrentCompany,
			Role:            user.Role,
			PostedJobsCount: int(jobCount), // 👈 Added this line
		}
	}

	var count int64
	config.DB.Model(&model.JobApplication{}).Where("user_id = ?", user.ID).Count(&count)

	return model.ApplicantResponse{
		ID:                user.ID,
		FullName:          user.FullName,
		Title:             user.Title,
		Location:          user.Location,
		Email:             user.Email,
		Phone:             user.Phone,
		CurrentCompany:    user.CurrentCompany,
		LinkedIn:          user.LinkedIn,
		GitHub:            user.GitHub,
		Portfolio:         user.Portfolio,
		Skills:            user.Skills,
		Image:             user.Image,
		Role:              user.Role,
		Education:         user.Education,
		Experience:        user.Experience,
		ApplicationsCount: int(count),
	}
}

// ---------- Signup ----------
func Signup(w http.ResponseWriter, r *http.Request) {
	var input SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Validate role
	allowedRoles := map[string]bool{"recruiter": true, "applicant": true, "admin": true}
	if _, ok := allowedRoles[input.Role]; !ok {
		http.Error(w, "Invalid role. Must be recruiter, applicant, or admin", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var existingUser model.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Password hashing failed", http.StatusInternalServerError)
		return
	}

	// Create user
	user := model.User{
		Email:          input.Email,
		Password:       string(hashedPassword),
		FullName:       input.FullName,
		Title:          input.Title,
		Location:       input.Location,
		Phone:          input.Phone,
		CurrentCompany: input.Company,
		LinkedIn:       input.LinkedIn,
		GitHub:         input.GitHub,
		Portfolio:      input.Portfolio,
		Skills:         input.Skills,
		Image:          input.Image,
		Role:           input.Role,
	}

	tx := config.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		return
	}

	// Save education
	for _, edu := range input.Education {
		edu.UserID = user.ID
		if err := tx.Create(&edu).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to add education", http.StatusInternalServerError)
			return
		}
	}

	// Save experience
	for _, exp := range input.Experience {
		exp.UserID = user.ID
		if err := tx.Create(&exp).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to add experience", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		return
	}

	// Reload with relations
	if err := config.DB.Preload("Education").Preload("Experience").First(&user, user.ID).Error; err != nil {
		http.Error(w, "Failed to load user after creation", http.StatusInternalServerError)
		return
	}

	// Generate JWT
	token := utils.GenerateJWT(user.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"token":   token,
		"user":    filterUserResponse(user),
	})
}

// ---------- Login ----------
func Login(w http.ResponseWriter, r *http.Request) {
	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token := utils.GenerateJWT(user.ID)
	if token == "" {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	// Load relations
	if err := config.DB.Preload("Education").Preload("Experience").First(&user).Error; err != nil {
		http.Error(w, "Failed to load user details", http.StatusInternalServerError)
		return
	}

	// Send filtered response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user":    filterUserResponse(user),
	})
}

// ---------- Reset Password ----------
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Password hashing failed", http.StatusInternalServerError)
		return
	}

	// Update password
	if err := config.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		http.Error(w, "Password update failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Password reset successfully",
		"email":   user.Email,
	})
}

// ---------- Get Profile ----------
// ---------- Get All Users ----------
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	if err := config.DB.Preload("Education").Preload("Experience").Find(&users).Error; err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	var filtered []interface{}
	for _, u := range users {
		filtered = append(filtered, filterUserResponse(u))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": filtered,
	})
}

// ---------- Get User By ID ----------
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path: /user/{id}
	path := r.URL.Path
	id := path[len("/user/"):]
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := config.DB.Preload("Education").Preload("Experience").First(&user, "id = ?", id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": filterUserResponse(user),
	})
}
