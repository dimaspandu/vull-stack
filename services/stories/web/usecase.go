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

func (u *Usecase) BrowseStories(ctx context.Context) (*StoryList, error) {
	result, err := u.Repo.List(ctx)
	if err != nil {
		return nil, err
	}

	data := StoryList{
		Stories: result.Data,
	}

	return &data, nil
}

func (u *Usecase) GetStoryBySlug(ctx context.Context, slug string) (*StoryShow, error) {
	result, err := u.Repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	var publishedAt string
	t, err := time.Parse(time.RFC3339, result.Data.PublishedAt)
	if err != nil {
		publishedAt = result.Data.PublishedAt
	} else {
		publishedAt = t.Format("January 2, 2006, 3:04 PM")
	}

	story := Story{
		ID:          result.Data.ID,
		Slug:        result.Data.Slug,
		Title:       result.Data.Title,
		Category:    result.Data.Category,
		Excerpt:     result.Data.Excerpt,
		PublishedAt: publishedAt,
		ReadingTime: result.Data.ReadingTime,
	}

	data := StoryShow{
		Story:   story,
		Content: result.Data.Content,
	}

	return &data, nil
}
