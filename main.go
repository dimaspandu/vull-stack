package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	api_stories "github.com/dimaspandu/vull-stack/services/stories/api"
	web_stories "github.com/dimaspandu/vull-stack/services/stories/web"
)

var (
	dbUser = "root"
	dbPass = ""
	dbName = "vss"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("client/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	dsn := fmt.Sprintf(
		"%s:%s@(127.0.0.1:3306)/%s?parseTime=true",
		dbUser,
		dbPass,
		dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	webStoriesRepo := web_stories.NewRepository(db)
	webStoriesUsecase := web_stories.NewUsecase(webStoriesRepo)
	webStoriesHandler := web_stories.NewHandler(webStoriesUsecase)

	mux.HandleFunc("GET /{$}", webStoriesHandler.List)
	mux.Handle("GET /favicon.ico", http.NotFoundHandler())
	mux.HandleFunc("GET /{slug}", webStoriesHandler.Show)

	apiStoriesRepo := api_stories.NewRepository(db)
	apiStoriesUseCase := api_stories.NewUsecase(apiStoriesRepo)
	apiStoriesHandler := api_stories.NewHandler(apiStoriesUseCase)

	mux.HandleFunc("GET /api/stories", apiStoriesHandler.List)
	mux.HandleFunc("GET /api/stories/{id}", apiStoriesHandler.Show)

	srv := http.Server{
		Addr:         ":6060",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Handler:      http.NewCrossOriginProtection().Handler(mux),
	}

	log.Fatal(srv.ListenAndServe())
}
