package web_stories

import (
	"context"
	"encoding/json"
	"os"
)

type Repository struct {
	Store string
}

func NewRepository(store string) *Repository {
	return &Repository{
		Store: store,
	}
}

func (r *Repository) List(ctx context.Context) (*StoryListResponse, error) {
	data, err := os.ReadFile(r.Store + ".json")
	if err != nil {
		return nil, err
	}

	var result StoryListResponse

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Repository) FindBySlug(ctx context.Context, slug string) (*StoryShowResponse, error) {
	data, err := os.ReadFile(r.Store + "/" + slug + ".json")
	if err != nil {
		return nil, err
	}

	var result StoryShowResponse

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
