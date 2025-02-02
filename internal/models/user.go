package models

import "time"

type User struct {
	ID        int       `json:"id,omitempty" db:"id"`
	Name      string    `json:"name,omitempty" db:"name"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}
