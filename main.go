package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

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

	sortStoryRepo := web_stories.NewRepository(db)
	sortStoryUsecase := web_stories.NewUsecase(sortStoryRepo)
	sortStoryHandler := web_stories.NewHandler(sortStoryUsecase)

	mux.HandleFunc("GET /{$}", sortStoryHandler.List)
	mux.Handle("GET /favicon.ico", http.NotFoundHandler())
	mux.HandleFunc("GET /{slug}", sortStoryHandler.Show)

	srv := http.Server{
		Addr:         ":6060",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Handler:      http.NewCrossOriginProtection().Handler(mux),
	}

	log.Fatal(srv.ListenAndServe())
}
