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

func (u *Usecase) BrowseStories(ctx context.Context) (*StoryListResponse, error) {
	stories, err := u.Repo.List(ctx)
	if err != nil {
		return nil, err
	}

	res := &StoryListResponse{
		Data: stories,
	}

	return res, nil
}

func (u *Usecase) GetStoryDetails(ctx context.Context, storyId string) (*StoryShowResponse, error) {
	story, err := u.Repo.FindBySlug(ctx, storyId)
	if err != nil {
		return nil, err
	}

	res := &StoryShowResponse{
		Data: story,
	}

	return res, nil
}
