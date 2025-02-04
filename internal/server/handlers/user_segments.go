package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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

var userSegmentsHandler = "user segments handler"

// NewUserSegmentsHandler creates a new UserSegmentsHandler instance.
func NewUserSegmentsHandler(ctx context.Context, uss userSegmentsService) *UserSegmentsHandler {
	return &UserSegmentsHandler{
		userSegments: uss,
		ctx:          ctx,
	}
}

// SegmentsRequest represents a request for updating user segments.
// @Description Segment lists for adding and deleting segments
type SegmentsRequest struct {
	// required: false
	Add []db.SegmentModification `json:"add,omitempty"`
	// required: false
	Remove []string `json:"remove,omitempty"`
}

// UpdateHandle processes user segment updates via HTTP request.
//
//	@Summary        Update user segments
//	@Description    Updates the user in the database and returns an instance of the user
//	@Tags           user-segments
//	@Accept         json
//	@Produce        json
//	@Param          id          path        int                 true    "User ID"
//	@Param          Segments    body        SegmentsRequest     true    "User change information"
//	@Success        200                                                 "User segments have been successfully changed"
//	@Router         /users/{id}/segments [patch]
func (uss *UserSegmentsHandler) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "UpdateHandle"

	var (
		err    error
		userID int
		sr     SegmentsRequest
	)
	if userID, err = strconv.Atoi(r.PathValue("id")); err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&sr); err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = uss.userSegments.Update(r.Context(), userID, sr.Add, sr.Remove); err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode("{ok:true}"); err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", userSegmentsHandler, "success", sr)
}

// GetActiveHandle retrieves active segments for a user via HTTP request.
//
//	@Summary        Get active user segments
//	@Description    Gets the active user segments by ID.
//	@Tags           user-segments
//	@Accept         json
//	@Produce        json
//	@Param          id      path        int                     true    "User ID"
//	@Success        200     {array}     dto.SegmentResponse             "Array with active user segments received"
//	@Router         /users/{id}/segments [get]
func (uss *UserSegmentsHandler) GetActiveHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "GetActiveHandle"

	var (
		err      error
		userID   int
		segments []*models.Segment
	)

	if userID, err = strconv.Atoi(r.PathValue("id")); err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if segments, err = uss.userSegments.GetActive(uss.ctx, userID); err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(segments); err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", userSegmentsHandler, "success", segments)
}

// GetHistoryCSVHandle generates a CSV report and returns JSON with the download URL.
//
//	@Summary        Update user segments
//	@Description    Updates the user in the database and returns an instance of the user
//	@Tags           user-segments-history
//	@Accept         json
//	@Produce        json
//	@Param          id      path        int     true    "User ID"
//	@Param          year    query       int     true    "Year, e.g. 2025"
//	@Param          month   query       int     true    "Month, e.g. 02"
//	@Success        200     {object}    dto.USHResponse "CSV-history is ready at the link"
//	@Router         /users/{id}/segments/history [get]
func (uss *UserSegmentsHandler) GetHistoryCSVHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "GetHistoryCSVHandle"

	var (
		err                 error
		userID, year, month int
		dURL                string
	)

	userID, err = strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the parameters of the request of the YEAR and Month
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")
	if yearStr == "" || monthStr == "" {
		slog.Error(fn, "handler", userSegmentsHandler, "err", "missing year or month")
		http.Error(w, "year and month parameters are required", http.StatusBadRequest)
		return
	}
	year, err = strconv.Atoi(yearStr)
	if err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}
	month, err = strconv.Atoi(monthStr)
	if err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, "invalid month", http.StatusBadRequest)
		return
	}

	// Calling the method of forming CSV reports.
	dURL, err = uss.userSegments.GetHistoryCSV(r.Context(), userID, year, month)
	if err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dURL = fmt.Sprintf("http://%s/%s", r.Host, dURL)

	// Return Json from URL to the CSV file.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]string{"url": dURL}); err != nil {
		slog.Error(fn, "handler", userSegmentsHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", userSegmentsHandler, "success", dURL)
}
