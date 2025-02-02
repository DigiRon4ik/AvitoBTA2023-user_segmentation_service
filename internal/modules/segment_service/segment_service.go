package segment_service

import (
	"context"

	"user_segmentation_service/internal/models"
)

type DB interface {
	CreateSegment(ctx context.Context, seg *models.Segment) error
	DeleteSegment(ctx context.Context, slug string) error
	UpdateSegment(ctx context.Context, seg *models.Segment) error
	GetSegmentBySlug(ctx context.Context, slug string) (*models.Segment, error)
	GetAllSegments(ctx context.Context) ([]*models.Segment, error)
}

type SegmentService struct {
	store DB
}

func NewSegmentService(store DB) *SegmentService {
	return &SegmentService{store: store}
}

func (s *SegmentService) Create(ctx context.Context, seg *models.Segment) (*models.Segment, error) {
	if err := s.store.CreateSegment(ctx, seg); err != nil {
		return nil, err
	}
	return seg, nil
}

func (s *SegmentService) Delete(ctx context.Context, slug string) error {
	return s.store.DeleteSegment(ctx, slug)
}

func (s *SegmentService) Update(ctx context.Context, seg *models.Segment) (*models.Segment, error) {
	if err := s.store.UpdateSegment(ctx, seg); err != nil {
		return nil, err
	}
	return seg, nil
}

func (s *SegmentService) GetByID(ctx context.Context, slug string) (*models.Segment, error) {
	return s.store.GetSegmentBySlug(ctx, slug)
}

func (s *SegmentService) GetAll(ctx context.Context) ([]*models.Segment, error) {
	return s.store.GetAllSegments(ctx)
}
