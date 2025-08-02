package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Experience struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	Title       string    `json:"title"`
	Years       string    `json:"years"`
	Description string    `json:"description"`
}

func (e *Experience) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	return
}
