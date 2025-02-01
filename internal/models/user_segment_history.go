package models

import "time"

type UserSegmentHistory struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	SegmentID int       `json:"segment_id" db:"segment_id"`
	Action    string    `json:"action" db:"action"` // "ADD" или "REMOVE"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
