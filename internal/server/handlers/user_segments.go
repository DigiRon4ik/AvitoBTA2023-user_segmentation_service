package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"user_segmentation_service/internal/db"
	"user_segmentation_service/internal/models"
)

type userSegmentsService interface {
	Update(ctx context.Context, userID int, add []db.SegmentModification, remove []string) error
	GetActive(ctx context.Context, userID int) ([]*models.Segment, error)
}

type UserSegmentsHandler struct {
	userSegments userSegmentsService
	ctx          context.Context
}

func NewUserSegmentsHandler(ctx context.Context, uss userSegmentsService) *UserSegmentsHandler {
	return &UserSegmentsHandler{
		userSegments: uss,
		ctx:          ctx,
	}
}

type SegmentsRequest struct {
	Add    []db.SegmentModification `json:"add,omitempty"`
	Remove []string                 `json:"remove,omitempty"`
}

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
