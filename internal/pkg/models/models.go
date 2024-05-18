package models

import "github.com/google/uuid"

type Musician struct {
	ID   uuid.UUID
	Name string
}

type Track struct {
	ID         uuid.UUID
	Title      string
	MusicianID uuid.UUID
	Played     int64
}
