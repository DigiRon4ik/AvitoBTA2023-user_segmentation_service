// Package handlers provide HTTP request handlers for user segments.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"user_segmentation_service/internal/models"
)

// segmentService defines the methods for interacting with the segment data.
type segmentService interface {
	Create(ctx context.Context, seg *models.Segment) error
	Delete(ctx context.Context, slug string) error
	Update(ctx context.Context, seg *models.Segment) error
	GetBySlug(ctx context.Context, slug string) (*models.Segment, error)
	GetAll(ctx context.Context) ([]*models.Segment, error)
}

// SegmentHandlers handles HTTP requests related to segments.
type SegmentHandlers struct {
	segments segmentService
	ctx      context.Context
}

// NewSegmentHandler initializes and returns a new SegmentHandlers instance.
func NewSegmentHandler(ctx context.Context, ss segmentService) *SegmentHandlers {
	return &SegmentHandlers{
		segments: ss,
		ctx:      ctx,
	}
}

// CreateHandle handles the request for creating a new segment.
// [ POST /segments ]
func (sh *SegmentHandlers) CreateHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		segment *models.Segment
	)

	if err = json.NewDecoder(r.Body).Decode(&segment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = sh.segments.Create(sh.ctx, segment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(segment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteHandle handles the request for deleting a segment.
// [ DELETE /segments/{slug} ]
func (sh *SegmentHandlers) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		slug = r.PathValue("slug")
	)

	if err = sh.segments.Delete(sh.ctx, slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// UpdateHandle handles the request for updating an existing segment.
// [ PUT /segments/{slug} ]
func (sh *SegmentHandlers) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		slug    = r.PathValue("slug")
		segment *models.Segment
	)

	if err = json.NewDecoder(r.Body).Decode(&segment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	segment.Slug = slug
	if err = sh.segments.Update(sh.ctx, segment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(segment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetHandle handles the request for retrieving a single segment by its slug.
// [ GET /segments/{slug} ]
func (sh *SegmentHandlers) GetHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		slug    = r.PathValue("slug")
		segment *models.Segment
	)

	if segment, err = sh.segments.GetBySlug(sh.ctx, slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(segment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAllHandle handles the request for retrieving all segments.
// [ GET /segments ]
func (sh *SegmentHandlers) GetAllHandle(w http.ResponseWriter, _ *http.Request) {
	var (
		err      error
		segments []*models.Segment
	)

	if segments, err = sh.segments.GetAll(sh.ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(segments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
