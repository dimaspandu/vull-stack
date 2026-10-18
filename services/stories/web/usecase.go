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

func (u *Usecase) ListStories(ctx context.Context, q string) (*StoryList, error) {
	stories, err := u.Repo.List(ctx, q)
	if err != nil {
		return nil, err
	}

	storyList := StoryList{
		Query:   q,
		Stories: stories,
	}

	return &storyList, nil
}

func (u *Usecase) GetStoryBySlug(ctx context.Context, slug string) (*Story, error) {
	story, err := u.Repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	var publishedAt string
	t, err := time.Parse(time.RFC3339, story.PublishedAt)
	if err != nil {
		publishedAt = story.PublishedAt
	} else {
		publishedAt = t.Format("January 2, 2006, 3:04 PM")
	}
	story.PublishedAt = publishedAt

	return story, nil
}
