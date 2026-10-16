package web_stories

import (
	"context"
	"html/template"
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
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data, err := h.Usecase.BrowseStories(ctx)
	if err != nil {
		http.Error(
			w,
			"Failed to load stories.",
			http.StatusInternalServerError,
		)
	}

	tmpl, err := template.ParseFiles("client/index.html")
	if err != nil {
		http.Error(
			w,
			"Failed to load template.",
			http.StatusInternalServerError,
		)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		return
	}
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	slug := r.PathValue("slug")

	data, err := h.Usecase.GetStoryBySlug(ctx, slug)
	if err != nil {
		http.Error(
			w,
			"Failed to load story.",
			http.StatusInternalServerError,
		)
		return
	}

	tmpl, err := template.ParseFiles("client/[id].html")
	if err != nil {
		http.Error(
			w,
			"Failed to load template.",
			http.StatusInternalServerError,
		)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		return
	}
}
