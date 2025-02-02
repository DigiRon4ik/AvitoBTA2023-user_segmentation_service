package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"user_segmentation_service/internal/models"
)

type userService interface {
	Create(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, userID int) error
	Update(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID int) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
}

type UserHandlers struct {
	users userService
	ctx   context.Context
}

func NewUserHandlers(ctx context.Context, us userService) *UserHandlers {
	return &UserHandlers{
		users: us,
		ctx:   ctx,
	}
}

func (uh *UserHandlers) CreateHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		user *models.User
	)

	if err = json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = uh.users.Create(uh.ctx, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (uh *UserHandlers) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		userID int
	)

	if userID, err = strconv.Atoi(r.PathValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = uh.users.Delete(uh.ctx, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandlers) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		userID int
		user   *models.User
	)

	if userID, err = strconv.Atoi(r.PathValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.ID = userID
	if err = uh.users.Update(uh.ctx, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandlers) GetHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		userID int
		user   *models.User
	)

	if userID, err = strconv.Atoi(r.PathValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user, err = uh.users.GetByID(uh.ctx, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandlers) GetAllHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		users []*models.User
	)

	if users, err = uh.users.GetAll(uh.ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
