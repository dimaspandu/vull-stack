package api_stories

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) List(ctx context.Context) ([]Story, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT
			id,
			slug,
			title,
			category,
			excerpt,
			published_at,
			reading_time
		FROM stories
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stories := make([]Story, 0)
	for rows.Next() {
		var story Story
		if err := rows.Scan(
			&story.ID,
			&story.Slug,
			&story.Title,
			&story.Category,
			&story.Excerpt,
			&story.PublishedAt,
			&story.ReadingTime,
		); err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stories, nil
}

func (r *Repository) GetByID(ctx context.Context, slug string) (*Story, error) {
	var story Story
	var contentJSON []byte

	if err := r.DB.QueryRowContext(ctx, `
		SELECT
			id,
			slug,
			title,
			category,
			excerpt,
			content,
			published_at,
			reading_time
		FROM stories
		WHERE slug = ?
	`, slug).Scan(
		&story.ID,
		&story.Slug,
		&story.Title,
		&story.Category,
		&story.Excerpt,
		&contentJSON,
		&story.PublishedAt,
		&story.ReadingTime,
	); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(contentJSON, &story.Content); err != nil {
		return nil, err
	}

	return &story, nil
}
