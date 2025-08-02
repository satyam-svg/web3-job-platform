package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Education struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	Institution string    `json:"institution"`
	Location    string    `json:"location"`
	Degree      string    `json:"degree"`
	GPA         string    `json:"gpa"`
	Years       string    `json:"years"`
}

func (e *Education) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	return
}
