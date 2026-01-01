package pet

import (
	"hona/backend/internal/domain/enums"
	"time"
)

type PetBasicDataResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Kind        string     `json:"kind"`
	Species     string     `json:"species"`
	Gender      string     `json:"gender"`
	PictureLink string     `json:"pictureLink"`
	BirthDate   *time.Time `json:"birthDate"`
	IsAdult     bool       `json:"isAdult"`
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

type PetKindResponse struct {
	Num  enums.PetKind `json:"num"`
	Name string        `json:"name"`
}

type PetSpeciesResponse struct {
	Num  enums.Species `json:"num"`
	Name string        `json:"name"`
}
