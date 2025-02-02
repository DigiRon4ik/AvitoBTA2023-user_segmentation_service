package models

import "time"

type Segment struct {
	ID          int       `json:"id,omitempty" db:"id"`
	Slug        string    `json:"slug,omitempty" db:"slug"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
}
