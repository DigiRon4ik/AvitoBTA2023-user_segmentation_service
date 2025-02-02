package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"user_segmentation_service/internal/models"
)

type segmentService interface {
	Create(ctx context.Context, seg *models.Segment) error
	Delete(ctx context.Context, slug string) error
	Update(ctx context.Context, seg *models.Segment) error
	GetBySlug(ctx context.Context, slug string) (*models.Segment, error)
	GetAll(ctx context.Context) ([]*models.Segment, error)
}

type SegmentHandlers struct {
	segments segmentService
	ctx      context.Context
}

func NewSegmentHandlers(ctx context.Context, ss segmentService) *SegmentHandlers {
	return &SegmentHandlers{
		segments: ss,
		ctx:      ctx,
	}
}

func (sh *SegmentHandlers) CreateHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		segment *models.Segment
	)

	if err = json.NewDecoder(r.Body).Decode(segment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = sh.segments.Create(sh.ctx, segment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(segment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (sh *SegmentHandlers) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		slug string = r.PathValue("slug")
	)

	if err = sh.segments.Delete(sh.ctx, slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (sh *SegmentHandlers) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		slug    string = r.PathValue("slug")
		segment *models.Segment
	)

	if err = json.NewDecoder(r.Body).Decode(segment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	segment.Slug = slug
	w.Header().Set("Content-Type", "application/json")
	if err = sh.segments.Update(sh.ctx, segment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (sh *SegmentHandlers) GetHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		slug    string = r.PathValue("slug")
		segment *models.Segment
	)

	if segment, err = sh.segments.GetBySlug(sh.ctx, slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(segment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (sh *SegmentHandlers) GetAllHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		segments []*models.Segment
	)

	if segments, err = sh.segments.GetAll(sh.ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(segments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
