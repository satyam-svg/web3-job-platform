package model

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	SalaryMin   int       `json:"salary_min"`
	SalaryMax   int       `json:"salary_max"`
	Type        string    `json:"type"` // Full-time, Part-time, etc.
	Description string    `json:"description"`
	Tags        string    `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`

	RecruiterID uuid.UUID `json:"recruiter_id"`                    // New field
	Recruiter   User      `gorm:"foreignKey:RecruiterID" json:"-"` // Avoid recursive json
}

type JobApplication struct {
	ID        uint      `gorm:"primaryKey"`
	JobID     uint      `json:"job_id"`
	UserID    uuid.UUID `json:"user_id"` // applicant
	CreatedAt time.Time

	Job  Job  `gorm:"foreignKey:JobID"`
	User User `gorm:"foreignKey:UserID"`
}
