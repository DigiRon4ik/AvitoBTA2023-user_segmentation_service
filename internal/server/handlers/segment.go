// Package handlers provide HTTP request handlers for user segments.
package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
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

var segmentHandler = "segment handler"

// NewSegmentHandler initializes and returns a new SegmentHandlers instance.
func NewSegmentHandler(ctx context.Context, ss segmentService) *SegmentHandlers {
	return &SegmentHandlers{
		segments: ss,
		ctx:      ctx,
	}
}

// CreateHandle handles the request for creating a new segment.
//
//	@Summary        Add segment
//	@Description    Creates a segment in the database and returns the instance
//	@Tags           segments
//	@Accept         json
//	@Produce        json
//	@Param          Segment body        dto.SegmentCreateRequest    true    "Information about the segment to be added"
//	@Success        201     {object}    dto.SegmentResponse                 "The segment has been successfully established"
//	@Router         /segments [post]
func (sh *SegmentHandlers) CreateHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "CreateHandle"

	var (
		err     error
		segment *models.Segment
	)

	if err = json.NewDecoder(r.Body).Decode(&segment); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = sh.segments.Create(sh.ctx, segment); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(segment); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", segmentHandler, "success", segment)
}

// DeleteHandle handles the request for deleting a segment.
//
//	@Summary        Delete segment
//	@Description    Deletes a segment from the database
//	@Tags           segments
//	@Accept         json
//	@Produce        json
//	@Param          slug    path    string  true    "Segment slug"
//	@Success        204                             "The segment with this slug has been successfully deleted"
//	@Router         /segments/{slug} [delete]
func (sh *SegmentHandlers) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "DeleteHandle"

	var (
		err  error
		slug = r.PathValue("slug")
	)

	if err = sh.segments.Delete(sh.ctx, slug); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	slog.Info(fn, "handler", segmentHandler, "success", slug)
}

// UpdateHandle handles the request for updating an existing segment.
//
//	@Summary        Update segment
//	@Description    Updates a segment in the database and returns an instance of it
//	@Tags           segments
//	@Accept         json
//	@Produce        json
//	@Param          slug    path        string                      true    "Segment slug"
//	@Param          Segment body        dto.SegmentUpdateRequest    true    "Segment change information"
//	@Success        200     {object}    dto.SegmentResponse                 "The segment with this slogan has been changed"
//	@Router         /segments/{slug} [put]
func (sh *SegmentHandlers) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "UpdateHandle"

	var (
		err     error
		slug    = r.PathValue("slug")
		segment *models.Segment
	)

	if err = json.NewDecoder(r.Body).Decode(&segment); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	segment.Slug = slug
	if err = sh.segments.Update(sh.ctx, segment); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(segment); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", segmentHandler, "success", segment)
}

// GetHandle handles the request for retrieving a single segment by its slug.
//
//	@Summary        Get segment
//	@Description    Get segment by slug
//	@Tags           segments
//	@Accept         json
//	@Produce        json
//	@Param          slug    path        string      true        "Segment slug"
//	@Success        200     {object}    dto.SegmentResponse     "A segment with such a slogan was obtained"
//	@Router         /segments/{slug} [get]
func (sh *SegmentHandlers) GetHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "GetHandle"

	var (
		err     error
		slug    = r.PathValue("slug")
		segment *models.Segment
	)

	if segment, err = sh.segments.GetBySlug(sh.ctx, slug); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(segment); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", segmentHandler, "success", segment)
}

// GetAllHandle handles the request for retrieving all segments.
//
//	@Summary        Get All segments
//	@Description    Get all segments from the database
//	@Tags           segments
//	@Accept         json
//	@Produce        json
//	@Success        200     {array}    dto.SegmentResponse     "An array of segments was obtained"
//	@Router         /segments [get]
func (sh *SegmentHandlers) GetAllHandle(w http.ResponseWriter, _ *http.Request) {
	const fn = "GetAllHandle"

	var (
		err      error
		segments []*models.Segment
	)

	if segments, err = sh.segments.GetAll(sh.ctx); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(segments); err != nil {
		slog.Error(fn, "handler", segmentHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", segmentHandler, "success", segments)
}
