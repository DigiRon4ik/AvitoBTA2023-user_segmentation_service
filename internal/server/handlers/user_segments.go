package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"user_segmentation_service/internal/db"
	"user_segmentation_service/internal/models"
)

// userSegmentsService defines methods for managing user segments.
type userSegmentsService interface {
	Update(ctx context.Context, userID int, add []db.SegmentModification, remove []string) error
	GetActive(ctx context.Context, userID int) ([]*models.Segment, error)
	GetHistoryCSV(ctx context.Context, userID, year, month int) (string, error)
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

// GetHistoryCSVHandle generates a CSV report and returns JSON with the download URL.
// [ GET /users/{id}/segments/history?year=YYYY&month=MM ]
func (uss *UserSegmentsHandler) GetHistoryCSVHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err                 error
		userID, year, month int
		dURL                string
	)

	userID, err = strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the parameters of the request of the YEAR and Month
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")
	if yearStr == "" || monthStr == "" {
		http.Error(w, "year and month parameters are required", http.StatusBadRequest)
		return
	}
	year, err = strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}
	month, err = strconv.Atoi(monthStr)
	if err != nil {
		http.Error(w, "invalid month", http.StatusBadRequest)
		return
	}

	// Calling the method of forming CSV reports.
	dURL, err = uss.userSegments.GetHistoryCSV(r.Context(), userID, year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dURL = fmt.Sprintf("http://%s/%s", r.Host, dURL)

	// Return Json from URL to the CSV file.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]string{"url": dURL}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
