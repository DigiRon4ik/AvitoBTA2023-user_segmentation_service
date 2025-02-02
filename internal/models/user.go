// Package models defines data structures for the application.
package models

import "time"

// User represents a user entity with basic information.
type User struct {
	ID        int       `json:"id,omitempty" db:"id"`
	Name      string    `json:"name,omitempty" db:"name"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}
