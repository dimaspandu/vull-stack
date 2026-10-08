package web_stories

import (
	"html/template"
	"net/http"
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data, err := h.Usecase.BrowseStories(r.Context())
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

	w.WriteHeader(http.StatusOK)

	if err := tmpl.Execute(w, data); err != nil {
		return
	}
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	slug := r.PathValue("slug")

	data, err := h.Usecase.GetStoryBySlug(r.Context(), slug)
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

	w.WriteHeader(http.StatusOK)

	if err := tmpl.Execute(w, data); err != nil {
		return
	}
}
