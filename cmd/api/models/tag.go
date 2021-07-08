package models

import "github.com/gofrs/uuid"

type Tag struct {
	ID  uuid.UUID `json:"ID"`
	Tag string    `json:"tag"`
}
