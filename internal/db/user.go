// Package db provides functionality for interacting with the PostgreSQL database.
package db

import (
	"context"

	"user_segmentation_service/internal/models"
)

const (
	createUser  = `INSERT INTO users (name) VALUES ($1) RETURNING id, created_at;`
	deleteUser  = `DELETE FROM users WHERE id = $1;`
	updateUser  = `UPDATE users SET name = $1 WHERE id = $2 RETURNING created_at;`
	getUserByID = `SELECT id, name, created_at FROM users WHERE id = $1;`
	getAllUsers = `SELECT * FROM users;`
)

// CreateUser creates a new user in the database.
// On successful execution, the ID and CreatedAt fields are populated into the user structure.
func (s *Store) CreateUser(ctx context.Context, user *models.User) error {
	return s.pool.QueryRow(ctx, createUser, user.Name).Scan(&user.ID, &user.CreatedAt)
}

// DeleteUser deletes a user by ID.
func (s *Store) DeleteUser(ctx context.Context, userID int) error {
	_, err := s.pool.Exec(ctx, deleteUser, userID)
	return err
}

// UpdateUser changes the user data (e.g. name) by id.
func (s *Store) UpdateUser(ctx context.Context, user *models.User) error {
	return s.pool.QueryRow(ctx, updateUser, user.Name, user.ID).Scan(&user.CreatedAt)
}

// GetUserByID returns the user by ID.
func (s *Store) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, getUserByID, userID).
		Scan(&user.ID, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllUsers returns all users.
func (s *Store) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := s.pool.Query(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0, 16)
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
