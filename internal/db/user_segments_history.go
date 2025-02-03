// Package db provides functionality for interacting with the PostgreSQL database.
package db

import (
	"context"
	"fmt"

	"user_segmentation_service/internal/models"
)

const (
	getUserSegmentHistory = `
		SELECT ush.user_id, u.name, s.slug, s.description, ush.action, ush.created_at
		FROM user_segments_history ush
		JOIN users u ON ush.user_id = u.id
		JOIN segments s ON ush.segment_id = s.id
		WHERE ush.user_id = $1 
		  AND EXTRACT(YEAR FROM ush.created_at) = $2
		  AND EXTRACT(MONTH FROM ush.created_at) = $3
		ORDER BY ush.created_at;`
)

// GetUserSegmentHistory receives a story for a given user and the period.
func (s *Store) GetUserSegmentHistory(ctx context.Context, userID, year, month int) ([]*models.HistoryRecord, error) {
	rows, err := s.pool.Query(ctx, getUserSegmentHistory, userID, year, month)
	if err != nil {
		return nil, fmt.Errorf("query history: %w", err)
	}
	defer rows.Close()

	records := make([]*models.HistoryRecord, 0, 32)
	for rows.Next() {
		rec := &models.HistoryRecord{}
		if err := rows.Scan(&rec.UserID, &rec.UserName, &rec.SegmentSlug, &rec.SegmentDescription, &rec.Action, &rec.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan history record: %w", err)
		}
		records = append(records, rec)
	}
	return records, nil
}
