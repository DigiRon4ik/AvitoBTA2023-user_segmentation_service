// Package models defines data structures for the application.
package models

import "time"

// UserSegmentHistory stores historical records of user-segment actions.
type UserSegmentHistory struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	SegmentID int       `json:"segment_id" db:"segment_id"`
	Action    string    `json:"action" db:"action"` // "ADD" или "REMOVE"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// HistoryRecord Describes a record for the formation of CSV reports.
type HistoryRecord struct {
	UserID             int       `json:"user_id"`
	UserName           string    `json:"user_name,omitempty"`
	SegmentSlug        string    `json:"segment_slug"`
	SegmentDescription string    `json:"segment_description,omitempty"`
	Action             string    `json:"action"`
	CreatedAt          time.Time `json:"created_at"`
}
