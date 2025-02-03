package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"user_segmentation_service/internal/db"
	"user_segmentation_service/internal/models"
)

// userSegmentsService defines methods for managing user segments.
type userSegmentsService interface {
	Update(ctx context.Context, userID int, add []db.SegmentModification, remove []string) error
	GetActive(ctx context.Context, userID int) ([]*models.Segment, error)
}

// UserSegmentsHandler handles HTTP requests for user segments.
type UserSegmentsHandler struct {
	userSegments userSegmentsService
	ctx          context.Context
}

// NewUserSegmentsHandler creates a new UserSegmentsHandler instance.
func NewUserSegmentsHandler(ctx context.Context, uss userSegmentsService) *UserSegmentsHandler {
	return &UserSegmentsHandler{
		userSegments: uss,
		ctx:          ctx,
	}
}

// SegmentsRequest represents a request for updating user segments.
type SegmentsRequest struct {
	Add    []db.SegmentModification `json:"add,omitempty"`
	Remove []string                 `json:"remove,omitempty"`
}

// UpdateHandle processes user segment updates via HTTP request.
// [ PATH /users{id}/segments ]
func (uss *UserSegmentsHandler) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		userID int
		sr     SegmentsRequest
	)
	if userID, err = strconv.Atoi(r.PathValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&sr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = uss.userSegments.Update(r.Context(), userID, sr.Add, sr.Remove); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode("{ok:true}"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetActiveHandle retrieves active segments for a user via HTTP request.
// [ GET /users{id}/segments ]
func (uss *UserSegmentsHandler) GetActiveHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		userID   int
		segments []*models.Segment
	)

	if userID, err = strconv.Atoi(r.PathValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if segments, err = uss.userSegments.GetActive(uss.ctx, userID); err != nil {
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
