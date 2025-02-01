package models

import "time"

type Segment struct {
	ID          int       `json:"id" db:"id"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
