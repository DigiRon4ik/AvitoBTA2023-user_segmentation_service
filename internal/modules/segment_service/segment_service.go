// Package segment_service provides business logic for managing user segments.
package segment_service

import (
	"context"

	"user_segmentation_service/internal/models"
)

// DB defines the required database operations for segment management.
type DB interface {
	CreateSegment(ctx context.Context, seg *models.Segment) error
	DeleteSegment(ctx context.Context, slug string) error
	UpdateSegment(ctx context.Context, seg *models.Segment) error
	GetSegmentBySlug(ctx context.Context, slug string) (*models.Segment, error)
	GetAllSegments(ctx context.Context) ([]*models.Segment, error)
}

// SegmentService handles operations related to user segments.
type SegmentService struct {
	store DB
}

// NewSegmentService creates a new instance of SegmentService.
func NewSegmentService(store DB) *SegmentService {
	return &SegmentService{store: store}
}

// Create adds a new segment to the database.
func (s *SegmentService) Create(ctx context.Context, seg *models.Segment) error {
	return s.store.CreateSegment(ctx, seg)
}

// Delete removes a segment by its slug.
func (s *SegmentService) Delete(ctx context.Context, slug string) error {
	return s.store.DeleteSegment(ctx, slug)
}

// Update modifies an existing segment.
func (s *SegmentService) Update(ctx context.Context, seg *models.Segment) error {
	return s.store.UpdateSegment(ctx, seg)
}

// GetBySlug retrieves a segment by its slug.
func (s *SegmentService) GetBySlug(ctx context.Context, slug string) (*models.Segment, error) {
	return s.store.GetSegmentBySlug(ctx, slug)
}

// GetAll returns all segments from the database.
func (s *SegmentService) GetAll(ctx context.Context) ([]*models.Segment, error) {
	return s.store.GetAllSegments(ctx)
}
