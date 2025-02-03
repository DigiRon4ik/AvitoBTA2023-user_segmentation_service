// Package user_segments_service provides business logic for managing user segments.
package user_segments_service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"user_segmentation_service/internal/db"
	"user_segmentation_service/internal/models"
)

// DB defines the required database operations for user management.
type DB interface {
	UpdateUserSegments(ctx context.Context, userID int, add []db.SegmentModification, remove []string) error
	GetActiveSegmentsForUser(ctx context.Context, userID int) ([]*models.Segment, error)
	GetUserSegmentHistory(ctx context.Context, userID, year, month int) ([]*models.HistoryRecord, error)
}

// UserSegmentationService encapsulates the business logic for handling user segmentation.
type UserSegmentationService struct {
	store DB
}

// NewUserSegmentationService creates a new service instance to handle user segmentation.
func NewUserSegmentationService(store DB) *UserSegmentationService {
	return &UserSegmentationService{
		store: store,
	}
}

// Update updates user segments by adding and removing segments.
// add - list of segments to add (with optional TTL),
// remove - list of slug segments to remove.
func (s *UserSegmentationService) Update(ctx context.Context, userID int, add []db.SegmentModification, remove []string) error {
	return s.store.UpdateUserSegments(ctx, userID, add, remove)
}

// GetActive returns the list of active user segments.
// Active segments are segments that have not yet expired (TTL).
func (s *UserSegmentationService) GetActive(ctx context.Context, userID int) ([]*models.Segment, error) {
	return s.store.GetActiveSegmentsForUser(ctx, userID)
}

// GetHistoryCSV generates a CSV report on the history of segment changes for the user
// for the specified year and month. The CSV file is saved in the "reports" directory, and the download URL is returned.
// TODO: Перенести в отдельный сервис.
func (s *UserSegmentationService) GetHistoryCSV(ctx context.Context, userID, year, month int) (string, error) {
	records, err := s.store.GetUserSegmentHistory(ctx, userID, year, month)
	if err != nil {
		return "", err
	}

	// Create a directorial for reports, if it is even not.
	outputDir := "reports"
	if err = os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", err
	}

	// Form the name of the file: report_{userID}_{year}_{month}.csv
	filename := fmt.Sprintf("report_%d_%d_%d.csv", userID, year, month)
	filePath := filepath.Join(outputDir, filename)

	var file *os.File
	file, err = os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	// Record the title
	header := []string{"user_id", "user_name", "segment_slug", "segment_description", "action", "created_at"}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	// Record the data
	for _, rec := range records {
		row := []string{
			strconv.Itoa(rec.UserID),
			rec.UserName,
			rec.SegmentSlug,
			rec.SegmentDescription,
			rec.Action,
			rec.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	// We assume that static files from the "reports" folder are served by the server at the URL /reports/
	downloadURL := fmt.Sprintf("reports/%s", filename)
	return downloadURL, nil
}
