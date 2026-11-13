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

	q := r.URL.Query().Get("q")

	result, err := h.Usecase.ListStories(ctx, q)
	if err != nil {
		http.Error(
			w,
			"Failed to load stories.",
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	slug := r.PathValue("slug")

	result, err := h.Usecase.GetStory(ctx, slug)
	if err != nil {
		http.Error(
			w,
			"Failed to load story.",
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(result)
}
