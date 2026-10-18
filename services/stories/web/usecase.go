package web_stories

import (
	"context"
	"time"
)

type Usecase struct {
	Repo *Repository
}

func NewUsecase(repo *Repository) *Usecase {
	return &Usecase{
		Repo: repo,
	}
}

func (u *Usecase) BrowseStories(ctx context.Context, q string) (*StoryList, error) {
	stories, err := u.Repo.List(ctx, q)
	if err != nil {
		return nil, err
	}

	data := StoryList{
		Query: q,
		Stories: stories,
	}

	return &data, nil
}

func (u *Usecase) GetStoryBySlug(ctx context.Context, slug string) (*Story, error) {
	data, err := u.Repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	var publishedAt string
	t, err := time.Parse(time.RFC3339, data.PublishedAt)
	if err != nil {
		publishedAt = data.PublishedAt
	} else {
		publishedAt = t.Format("January 2, 2006, 3:04 PM")
	}
	data.PublishedAt = publishedAt

	return data, nil
}
