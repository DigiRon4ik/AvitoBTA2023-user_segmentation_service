package models

import "time"

type UserSegment struct {
	UserID         int       `json:"user_id" db:"user_id"`
	SegmentID      int       `json:"segment_id" db:"segment_id"`
	ExpirationTime time.Time `json:"expiration_time" db:"expiration_time"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
