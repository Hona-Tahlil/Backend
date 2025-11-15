package pet

import (
	"hona/backend/internal/domain/enums"
	"time"
)

type PetBasicDataResponse struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Kind        enums.PetKind   `json:"kind"`
	Species     enums.Species   `json:"species"`
	Gender      enums.PetGender `json:"gender"`
	PictureLink string          `json:"pictureLink"`
}

type PetFullDataResponse struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Kind        enums.PetKind   `json:"kind"`
	Species     enums.Species   `json:"species"`
	Gender      enums.PetGender `json:"gender"`
	PictureLink string          `json:"pictureLink"`
	BirthDate   *time.Time      `json:"birthDate"`
	IsAdult     bool            `json:"isAdult"`
	Weight      *float32        `json:"weight"`
	AboutPet    *string         `json:"aboutPet"`
}
