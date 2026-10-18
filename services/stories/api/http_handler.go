package api_stories

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	Usecase *Usecase
}

func NewHandler(usecase *Usecase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	res, err := h.Usecase.BrowseStories(ctx)
	if err != nil {
		http.Error(
			w,
			"Failed to load stories.",
			http.StatusInternalServerError,
		)
	}

	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	storyId := r.PathValue("id")

	res, err := h.Usecase.GetStoryDetails(ctx, storyId)
	if err != nil {
		http.Error(
			w,
			"Failed to load story.",
			http.StatusInternalServerError,
		)
	}

	json.NewEncoder(w).Encode(res)
}
