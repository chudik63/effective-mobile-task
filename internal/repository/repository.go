package repository

import (
	"context"
	"database/sql"
	"effective-mobile-task/internal/database/postgres"
	"effective-mobile-task/internal/models"
	"errors"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type Creds map[string]interface{}

type Repository struct {
	db postgres.DB
}

func New(db postgres.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) AddSong(ctx context.Context, song *models.Song) (uint64, error) {
	var id uint64
	err := sq.Insert("public.songs").
		Columns("group_name", "song", "release_date", "text", "link").
		Values(song.Group, song.Song, song.ReleaseDate, song.Text, song.Link).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow().
		Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateSong(ctx context.Context, song *models.Song) error {
	query := sq.Update("public.songs").PlaceholderFormat(sq.Dollar)

	if song.Group != "" {
		query = query.Set("group_name", song.Group)
	}
	if song.Song != "" {
		query = query.Set("song", song.Song)
	}
	if !song.ReleaseDate.IsZero() {
		query = query.Set("release_date", song.ReleaseDate)
	}
	if song.Text != "" {
		query = query.Set("text", song.Text)
	}
	if song.Link != "" {
		query = query.Set("link", song.Link)
	}

	query = query.Where(sq.Eq{"id": song.Id})

	res, err := query.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *Repository) DeleteSong(ctx context.Context, id uint64) error {
	res, err := sq.Delete("public.songs").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *Repository) GetSongText(ctx context.Context, id uint64) (string, error) {
	var text string

	err := sq.Select("text").
		From("public.songs").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow().
		Scan(&text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNotFound
		}

		return "", err
	}

	return text, nil
}

func (r *Repository) GetSongs(ctx context.Context, creds Creds, offset uint64, limit uint64) ([]*models.Song, error) {
	query := sq.Select("*").
		From("public.songs").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	for key, value := range creds {
		query = query.Where(sq.Like{key: "%" + value.(string) + "%"})
	}

	query = query.Limit(limit).Offset(offset)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []*models.Song

	for rows.Next() {
		var song models.Song

		if err := rows.Scan(&song.Id, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, err
		}

		songs = append(songs, &song)
	}

	if rows.Err() != nil {
		return nil, err
	}

	if len(songs) == 0 {
		return nil, models.ErrNotFound
	}

	return songs, nil
}
