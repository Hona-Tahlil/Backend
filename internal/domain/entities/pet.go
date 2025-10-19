package entities

import (
	"time"

	"gorm.io/gorm"
)

type Pet struct {
	gorm.Model
	UserID       uint
	Name         string
	Kind         string // make it enum later
	Species      string
	BirthDate    time.Time
	Gender       string // make enum later
	Weight       int
	HealthRecord string
	PictureLink  string
}
