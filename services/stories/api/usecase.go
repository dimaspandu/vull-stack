package api_stories

import "context"

type Usecase struct {
	Repo *Repository
}

func NewUsecase(repo *Repository) *Usecase {
	return &Usecase{
		Repo: repo,
	}
}

func (u *Usecase) ListStories(ctx context.Context, q string) (*StoryListResponse, error) {
	stories, err := u.Repo.List(ctx, q)
	if err != nil {
		return nil, err
	}

	out := &StoryListResponse{
		Data: stories,
	}

	return out, nil
}

func (u *Usecase) GetStory(ctx context.Context, slug string) (*StoryShowResponse, error) {
	story, err := u.Repo.GetByID(ctx, slug)
	if err != nil {
		return nil, err
	}

	out := &StoryShowResponse{
		Data: story,
	}

	return out, nil
}
