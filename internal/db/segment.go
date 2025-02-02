// Package db provides functionality for interacting with the PostgreSQL database.
package db

import (
	"bytes"
	"context"
	"io"

	"user_segmentation_service/internal/models"
)

const (
	createSegment    = `INSERT INTO segments (slug, description) VALUES ($1, $2) RETURNING id, created_at;`
	deleteSegment    = `DELETE FROM segments WHERE slug = $1;`
	updateSegment    = `UPDATE segments SET description = $1 WHERE slug = $2 RETURNING id, created_at;`
	getSegmentBySlug = `SELECT id, slug, description, created_at FROM segments WHERE slug = $1;`
	getAllSegments   = `SELECT * FROM segments;`
)

// CreateSegment creates a new segment in the database.
// On successful execution, the ID and CreatedAt fields are populated into the seg structure.
func (s *Store) CreateSegment(ctx context.Context, seg *models.Segment) error {
	return s.pool.QueryRow(ctx, createSegment, seg.Slug, seg.Description).Scan(&seg.ID, &seg.CreatedAt)
}

// DeleteSegment deletes a segment from the database by slug.
func (s *Store) DeleteSegment(ctx context.Context, slug string) error {
	_, err := s.pool.Exec(ctx, deleteSegment, slug)
	return err
}

// UpdateSegment changes the segment data (e.g., description) by slug.
// Here only the description field is updated, but others can be added if necessary.
func (s *Store) UpdateSegment(ctx context.Context, seg *models.Segment) error {
	return s.pool.QueryRow(ctx, updateSegment, seg.Description, seg.Slug).Scan(&seg.ID, &seg.CreatedAt)
}

// GetSegmentBySlug gets the segment by slug.
func (s *Store) GetSegmentBySlug(ctx context.Context, slug string) (*models.Segment, error) {
	seg := &models.Segment{}
	err := s.pool.QueryRow(ctx, getSegmentBySlug, slug).
		Scan(&seg.ID, &seg.Slug, &seg.Description, &seg.CreatedAt)
	if err != nil {
		return nil, err
	}
	return seg, nil
}

// GetAllSegments returns a list of all segments.
func (s *Store) GetAllSegments(ctx context.Context) ([]*models.Segment, error) {
	rows, err := s.pool.Query(ctx, getAllSegments)
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

// GetAllSegmentsViaCopy retrieves all segments using COPY.
func (s *Store) GetAllSegmentsViaCopy(ctx context.Context) (io.Reader, error) {
	// Form a SQL request for Copy.
	// Use row_to_json to get each row as JSON.
	query := `COPY (SELECT row_to_json(s) FROM (SELECT id, slug, description, created_at FROM segments) s) TO STDOUT`

	// Buffer for data retrieval.
	var buf bytes.Buffer

	// Execute the COPY command and transfer the result to the buffer.
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	_, err = conn.Conn().PgConn().CopyTo(ctx, &buf, query)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
