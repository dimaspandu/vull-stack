package web_stories

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
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

	q := r.URL.Query().Get("q")

	pageData, err := h.Usecase.ListStories(ctx, q)
	if err != nil {
		h.InternalServerError(w, "Failed to load stories.")
		return
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

	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(
			w,
			"Failed to render template.",
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	slug := r.PathValue("slug")

	story, err := h.Usecase.GetStoryBySlug(ctx, slug)
	if err != nil {
		// An unknown slug is a missing page, not a server failure.
		if errors.Is(err, sql.ErrNoRows) {
			h.NotFound(w, r)
			return
		}

		h.InternalServerError(w, "Failed to load story.")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("client/[slug].html")
	if err != nil {
		http.Error(
			w,
			"Failed to load template.",
			http.StatusInternalServerError,
		)
		return
	}

	if err := tmpl.Execute(w, story); err != nil {
		http.Error(
			w,
			"Failed to render template.",
			http.StatusInternalServerError,
		)
		return
	}
}

// NotFound renders the static 404 page. It is also the fallback route, so any
// unmatched path or method ends up here.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("client/404.html")
	if err != nil {
		http.Error(
			w,
			"Failed to load template.",
			http.StatusInternalServerError,
		)
		return
	}

	// Render into a buffer first: the status code must not be sent before the
	// page is known to be complete.
	var page bytes.Buffer
	if err := tmpl.Execute(&page, nil); err != nil {
		http.Error(
			w,
			"Failed to render template.",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(page.Bytes())
}

// InternalServerError renders the static 500 page for server failures.
func (h *Handler) InternalServerError(w http.ResponseWriter, message string) {
	tmpl, err := template.ParseFiles("client/500.html")
	if err != nil {
		http.Error(
			w,
			"Failed to load template.",
			http.StatusInternalServerError,
		)
		return
	}

	var page bytes.Buffer
	if err := tmpl.Execute(&page, nil); err != nil {
		http.Error(
			w,
			"Failed to render template.",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(page.Bytes())
}