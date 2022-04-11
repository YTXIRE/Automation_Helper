package store

import (
	"backend/internal/app/models"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *models.User) (*models.User, error) {
	if err := r.store.db.QueryRow(
		"INSERT INTO public.users (email, password, login, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.Password,
		u.Login,
		u.CreatedAt,
	).Scan(&u.ID); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, login, created_at FROM public.users WHERE email = $1",
		email,
	).Scan(&u.ID, &u.Email, &u.Login, &u.CreatedAt); err != nil {
		return nil, err
	}
	return u, nil
}
