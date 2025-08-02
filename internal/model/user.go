package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	FullName       string    `json:"full_name"`
	Title          string    `json:"title"`
	Location       string    `json:"location"`
	Email          string    `gorm:"unique" json:"email"`
	Password       string    `json:"-"` // Hashed password
	Phone          string    `json:"phone"`
	CurrentCompany string    `json:"current_company"`
	LinkedIn       string    `json:"linkedin"`
	GitHub         string    `json:"github"`
	Portfolio      string    `json:"portfolio"`
	Skills         string    `json:"skills"`
	Image          string    `json:"image"`
	Role           string    `json:"role"`
	Credits        int       `json:"credits" gorm:"default:5"` // 👈 New field

	Education  []Education  `json:"education" gorm:"foreignKey:UserID"`
	Experience []Experience `json:"experience" gorm:"foreignKey:UserID"`
}

type RecruiterResponse struct {
	ID              uuid.UUID `json:"id"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Image           string    `json:"image"`
	CurrentCompany  string    `json:"current_company"`
	Role            string    `json:"role"`
	PostedJobsCount int       `json:"posted_jobs_count"` // 👈 ADD THIS
}

// ApplicantResponse struct for applicants/admin
type ApplicantResponse struct {
	ID                uuid.UUID    `json:"id"`
	FullName          string       `json:"full_name"`
	Title             string       `json:"title"`
	Location          string       `json:"location"`
	Email             string       `json:"email"`
	Phone             string       `json:"phone"`
	CurrentCompany    string       `json:"current_company"`
	LinkedIn          string       `json:"linkedin"`
	GitHub            string       `json:"github"`
	Portfolio         string       `json:"portfolio"`
	Skills            string       `json:"skills"`
	Image             string       `json:"image"`
	Role              string       `json:"role"`
	Education         []Education  `json:"education"`
	Experience        []Experience `json:"experience"`
	ApplicationsCount int          `json:"applications_count"` // new
}

// Automatically generate UUID before creating
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	if u.Credits == 0 {
		u.Credits = 5 // Only if it's not already set
	}
	return
}
