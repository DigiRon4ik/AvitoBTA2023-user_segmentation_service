// Package db provides functionality for interacting with the PostgreSQL database.
package db

import (
	"context"
	"fmt"
	"time"

	"user_segmentation_service/internal/models"
)

const (
	getActiveSegmentsForUser = `
		SELECT s.id, s.slug, s.description, s.created_at
		FROM segments s
		JOIN user_segments us ON s.id = us.segment_id
		WHERE us.user_id = $1 AND us.expiration_time > NOW()`
	// Удаляет записи из user_segments для заданного user_id и списка slug'ов,
	// возвращая удалённые данные (user_id, segment_id, created_at).
	// Затем сразу же записывает эти данные в user_segments_history с пометкой 'REMOVE'.
	// Используем CTE (WITH deleted_segments) для объединения удаления и логирования в один запрос.
	removingSegmentsForUser = `
		WITH deleted_segments AS (
            DELETE FROM user_segments
            WHERE user_id = $1
				AND segment_id IN (SELECT id
									FROM segments
									WHERE slug = ANY ($2))
            RETURNING user_id, segment_id, created_at)
		INSERT INTO user_segments_history (user_id, segment_id, action, created_at)
		SELECT user_id, segment_id, 'REMOVE', created_at
		FROM deleted_segments;`
	// Массовое добавление или обновление записей в user_segments с записью в историю.
	// 1. Преобразуем массивы slug и expiration_time в таблицу (segments_data).
	// 2. Находим segment_id по slug'ам (segment_ids).
	// 3. Вставляем новые или обновляем существующие записи в user_segments (inserted_segments).
	// 4. Фиксируем успешные операции в user_segments_history.
	addingSegmentsForUser = `
		WITH segments_data AS (SELECT UNNEST($1::TEXT[]) AS slug,
									UNNEST($2::TIMESTAMP[]) AS expiration_time),
			segment_ids AS (SELECT sd.slug,
								sd.expiration_time,
								s.id AS segment_id
							FROM segments_data sd
								JOIN segments s ON sd.slug = s.slug),
			inserted_segments AS (
				INSERT INTO user_segments (user_id, segment_id, expiration_time)
				SELECT $3 AS user_id, si.segment_id, si.expiration_time
				FROM segment_ids si
				ON CONFLICT (user_id, segment_id)
				DO UPDATE SET expiration_time = excluded.expiration_time
                RETURNING user_id, segment_id, created_at)
		INSERT INTO user_segments_history (user_id, segment_id, action, created_at)
		SELECT user_id, segment_id, 'ADD' AS action, created_at
		FROM inserted_segments i
		WHERE NOT EXISTS (
            SELECT 1 FROM user_segments_history h
            WHERE h.user_id = i.user_id
                AND h.segment_id = i.segment_id
                AND h.action = 'ADD'
                AND h.created_at = i.created_at
		);`
)

// SegmentModification describes the data for adding a segment to a user.
type SegmentModification struct {
	// required: true
	Slug string `json:"slug"` // segment slug
	// required: false
	ExpirationTime *time.Time `json:"expiration_time,omitempty"` // Optionally, if nil, the default value is used
}

// defaultExpiration specifies a default expiration time of 100 years from the current point in time.
func defaultExpiration() time.Time {
	return time.Now().Add(100 * 365 * 24 * time.Hour) // Approximately 100 years
}

// UpdateUserSegments updates user segments (transaction): adds and deletes segments.
// For each added segment, a record is inserted into the user_segments table and recorded in the history.
// For each segment to be deleted, the connection is deleted and the deletion is recorded in the history.
// TODO: Разделить на маленькие функции.
func (s *Store) UpdateUserSegments(ctx context.Context, userID int, add []SegmentModification, remove []string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("the beginning of the transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	// TODO: добавьте использование batch-запросов или COPY FROM если len(add/remove) > 10k+
	// Removing segments
	_, err = tx.Exec(ctx, removingSegmentsForUser, userID, remove)
	if err != nil {
		return fmt.Errorf("error delete segments for user %d: %w", userID, err)
	}

	// Adding segments
	if len(add) == 0 {
		return nil
	}
	// Data preparation for request
	slugs := make([]string, len(add))
	expTimes := make([]*time.Time, len(add))
	for i, mod := range add {
		slugs[i] = mod.Slug
		if mod.ExpirationTime == nil {
			defaultExp := defaultExpiration()
			expTimes[i] = &defaultExp
		} else {
			expTimes[i] = mod.ExpirationTime
		}
	}
	// Request
	_, err = tx.Exec(ctx, addingSegmentsForUser, slugs, expTimes, userID)
	if err != nil {
		return fmt.Errorf("error adding segments for user %d: %w", userID, err)
	}
	return nil
}

// GetActiveSegmentsForUser returns active user segments.
// Segments with an expiration time greater than the current time are considered active.
func (s *Store) GetActiveSegmentsForUser(ctx context.Context, userID int) ([]*models.Segment, error) {
	rows, err := s.pool.Query(ctx, getActiveSegmentsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	segments := make([]*models.Segment, 0, 16)
	for rows.Next() {
		seg := &models.Segment{}
		if err := rows.Scan(&seg.ID, &seg.Slug, &seg.Description, &seg.CreatedAt); err != nil {
			return nil, err
		}
		segments = append(segments, seg)
	}
	return segments, nil
}
