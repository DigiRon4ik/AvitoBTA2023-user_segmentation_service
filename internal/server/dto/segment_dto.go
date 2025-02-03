// Package dto for Swagger
package dto

import "time"

// SegmentCreateRequest for Swagger
//
//	@Description Segment information at creation
type SegmentCreateRequest struct {
	// required: true
	Slug string `json:"slug"`
	// required: false
	Description string `json:"description,omitempty"`
}

// SegmentUpdateRequest for Swagger
//
//	@Description Segment information when updating
type SegmentUpdateRequest struct {
	// required: true
	Description string `json:"description"`
}

// SegmentResponse for Swagger
//
//	@Description Segment information when creating/updating a segment
type SegmentResponse struct {
	// read only: true
	ID          int    `json:"id"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	// read only: true
	CreatedAt time.Time `json:"created_at"`
}
