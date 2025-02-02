package handlers

import (
	"context"
	"net/http"

	"user_segmentation_service/internal/models"
)

type segmentService interface {
	Create(ctx context.Context, seg *models.Segment) (*models.Segment, error)
	Delete(ctx context.Context, slug string) error
	Update(ctx context.Context, seg *models.Segment) (*models.Segment, error)
	GetByID(ctx context.Context, slug string) (*models.Segment, error)
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

func (uh *SegmentHandlers) CreateHandle(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *SegmentHandlers) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *SegmentHandlers) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *SegmentHandlers) GetHandle(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *SegmentHandlers) GetAllHandle(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
