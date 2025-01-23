package service

import (
	"context"
	"effective-mobile-task/internal/models"
	"effective-mobile-task/internal/repository"
	"strings"
)

type Repository interface {
	AddSong(ctx context.Context, song *models.Song) (uint64, error)
	UpdateSong(ctx context.Context, song *models.Song) error
	DeleteSong(ctx context.Context, id uint64) error
	GetSongText(ctx context.Context, id uint64) (string, error)
	GetSongs(ctx context.Context, creds repository.Creds, offset uint64, limit uint64) ([]*models.Song, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetSongLyrics(ctx context.Context, id uint64, offset int) (string, error) {
	lyrics, err := s.repo.GetSongText(ctx, id)
	if err != nil {
		return "", err
	}

	verses := strings.Split(lyrics, "\n\n")

	return verses[offset], nil
}

func (s *Service) UpdateSong(ctx context.Context, song *models.Song) error {
	return s.repo.UpdateSong(ctx, song)
}

func (s *Service) AddSong(ctx context.Context, song *models.Song) (uint64, error) {
	return s.repo.AddSong(ctx, song)
}

func (s *Service) DeleteSong(ctx context.Context, id uint64) error {
	return s.repo.DeleteSong(ctx, id)
}

func (s *Service) GetSongs(ctx context.Context, creds repository.Creds, offset uint64, limit uint64) ([]*models.Song, error) {
	return s.repo.GetSongs(ctx, creds, offset, limit)
}
