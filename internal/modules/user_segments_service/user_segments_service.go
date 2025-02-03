// Package user_segments_service provides business logic for managing user segments.
package user_segments_service

import (
	"context"

	"user_segmentation_service/internal/db"
	"user_segmentation_service/internal/models"
)

// DB defines the required database operations for user management.
type DB interface {
	UpdateUserSegments(ctx context.Context, userID int, add []db.SegmentModification, remove []string) error
	GetActiveSegmentsForUser(ctx context.Context, userID int) ([]*models.Segment, error)
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

// GetUserSegmentsHistory возвращает историю изменений сегментов для пользователя за указанный год и месяц.
// Метод предполагает, что в слое DB реализован соответствующий функционал.
// func (s *UserSegmentationService) GetUserSegmentsHistory(ctx context.Context, userID int, year int, month int) ([]*models.UserSegmentHistory, error) {
// 	return s.store.GetUserSegmentHistory(ctx, userID, year, month)
// }
