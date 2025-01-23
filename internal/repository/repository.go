package repository

import "effective-mobile-task/internal/database/postgres"

type Repository struct {
	db postgres.DB
}

func New(db postgres.DB) *Repository {
	return &Repository{db}
}
