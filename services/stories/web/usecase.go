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

	for i := range stories {
		stories[i].PublishedAt = formatPublishedAt(stories[i].PublishedAt)
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

	story.PublishedAt = formatPublishedAt(story.PublishedAt)

	return story, nil
}

func formatPublishedAt(publishedAt string) string {
	t, err := time.Parse(time.RFC3339, publishedAt)
	if err != nil {
		return publishedAt
	}
	return t.Format("January 2, 2006, 3:04 PM")
}
