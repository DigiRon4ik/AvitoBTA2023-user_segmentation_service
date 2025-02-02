// Package models defines data structures for the application.
package models

import "time"

// Segment represents a user segment with metadata.
type Segment struct {
	ID          int       `json:"id,omitempty" db:"id"`
	Slug        string    `json:"slug,omitempty" db:"slug"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
}
