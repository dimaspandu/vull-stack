package main

import (
	"log"
	"net/http"
	"time"

	sort_story "github.com/dimaspandu/vull-stack/services/stories/web"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("client/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	sortStoryRepo := sort_story.NewRepository("data/stories")
	sortStoryUsecase := sort_story.NewUsecase(sortStoryRepo)
	sortStoryHandler := sort_story.NewHandler(sortStoryUsecase)

	mux.HandleFunc("/", sortStoryHandler.List)
	mux.HandleFunc("/{slug}", sortStoryHandler.Show)

	srv := http.Server{
		Addr:         ":6060",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Handler:      http.NewCrossOriginProtection().Handler(mux),
	}

	log.Fatal(srv.ListenAndServe())
}
