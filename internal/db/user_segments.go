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
)

// segmentModification describes the data for adding a segment to a user.
type segmentModification struct {
	Slug           string     // segment slug
	ExpirationTime *time.Time // Optionally, if nil, the default value is used
}

// defaultExpiration specifies a default expiration time of 100 years from the current point in time.
func defaultExpiration() time.Time {
	return time.Now().Add(100 * 365 * 24 * time.Hour) // Approximately 100 years
}

// UpdateUserSegments updates user segments (transaction): adds and deletes segments.
// For each added segment, a record is inserted into the user_segments table and recorded in the history.
// For each segment to be deleted, the connection is deleted and the deletion is recorded in the history.
func (s *Store) UpdateUserSegments(ctx context.Context, userID int, add []segmentModification, remove []string) error {
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

	// Removing segments
	for _, slug := range remove {
		// Get segment id by slug
		var segmentID int
		err = tx.QueryRow(ctx, `SELECT id FROM segments WHERE slug = $1`, slug).Scan(&segmentID)
		if err != nil {
			return fmt.Errorf("no segment found with slug %s: %w", slug, err)
		}

		// Deleting the link
		_, err = tx.Exec(ctx, `
			DELETE FROM user_segments
			WHERE user_id = $1 AND segment_id = $2
		`, userID, segmentID)
		if err != nil {
			return fmt.Errorf("delete segment %s for user %d: %w", slug, userID, err)
		}

		// Record the operation in the history
		_, err = tx.Exec(ctx, `
			INSERT INTO user_segments_history (user_id, segment_id, action)
			VALUES ($1, $2, 'REMOVE')
		`, userID, segmentID)
		if err != nil {
			return fmt.Errorf("recording the history of deletion of segment %s: %w", slug, err)
		}
	}

	// Adding segments
	for _, mod := range add {
		// Get segment id by slug
		var segmentID int
		err = tx.QueryRow(ctx, `SELECT id FROM segments WHERE slug = $1`, mod.Slug).Scan(&segmentID)
		if err != nil {
			return fmt.Errorf("no segment found with slug %s: %w", mod.Slug, err)
		}

		// If expiration_time is not specified, we use the default value.
		expTime := mod.ExpirationTime
		if expTime == nil {
			defaultExp := defaultExpiration()
			expTime = &defaultExp
		}

		// Insert or update a link in the user_segments table.
		// If the record already exists, you can update expiration_time.
		_, err = tx.Exec(ctx, `
			INSERT INTO user_segments (user_id, segment_id, expiration_time)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, segment_id) DO UPDATE SET expiration_time = excluded.expiration_time
		`, userID, segmentID, *expTime)
		if err != nil {
			return fmt.Errorf("adding segment %s for user %d: %w", mod.Slug, userID, err)
		}

		// Record the operation in the history
		_, err = tx.Exec(ctx, `
			INSERT INTO user_segments_history (user_id, segment_id, action)
			VALUES ($1, $2, 'ADD')
		`, userID, segmentID)
		if err != nil {
			return fmt.Errorf("recording the history of adding segment %s: %w", mod.Slug, err)
		}
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
