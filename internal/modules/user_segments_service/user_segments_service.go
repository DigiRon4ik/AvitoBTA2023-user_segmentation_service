package user_segments_service

import (
	"context"

	"user_segmentation_service/internal/db"
	"user_segmentation_service/internal/models"
)

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

// Update обновляет сегменты пользователя, выполняя добавление и удаление сегментов.
// add - список сегментов для добавления (с опциональным TTL),
// remove - список slug сегментов, которые нужно удалить.
func (s *UserSegmentationService) Update(ctx context.Context, userID int, add []db.SegmentModification, remove []string) error {
	return s.store.UpdateUserSegments(ctx, userID, add, remove)
}

// GetActive возвращает список активных сегментов пользователя.
// Активными считаются сегменты, у которых время истечения (TTL) ещё не наступило.
func (s *UserSegmentationService) GetActive(ctx context.Context, userID int) ([]*models.Segment, error) {
	return s.store.GetActiveSegmentsForUser(ctx, userID)
}

// GetUserSegmentsHistory возвращает историю изменений сегментов для пользователя за указанный год и месяц.
// Метод предполагает, что в слое DB реализован соответствующий функционал.
// func (s *UserSegmentationService) GetUserSegmentsHistory(ctx context.Context, userID int, year int, month int) ([]*models.UserSegmentHistory, error) {
// 	return s.store.GetUserSegmentHistory(ctx, userID, year, month)
// }
