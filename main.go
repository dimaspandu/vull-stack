package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	api_stories "github.com/dimaspandu/vull-stack/services/stories/api"
	web_stories "github.com/dimaspandu/vull-stack/services/stories/web"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	_ = godotenv.Load(".env")
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("client/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		getenv("DB_USER", "root"),
		getenv("DB_PASS", ""),
		getenv("DB_HOST", "127.0.0.1"),
		getenv("DB_PORT", "3306"),
		getenv("DB_NAME", "vss"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Web stories
	webRepo := web_stories.NewRepository(db)
	webUC := web_stories.NewUsecase(webRepo)
	webH := web_stories.NewHandler(webUC)

	mux.HandleFunc("GET /{$}", webH.List)
	mux.Handle("GET /favicon.ico", http.NotFoundHandler())
	mux.HandleFunc("GET /{slug}", webH.Show)

	// API stories
	apiRepo := api_stories.NewRepository(db)
	apiUC := api_stories.NewUsecase(apiRepo)
	apiH := api_stories.NewHandler(apiUC)

	mux.HandleFunc("GET /api/stories", apiH.List)
	mux.HandleFunc("GET /api/stories/{id}", apiH.Show)

	port := getenv("APP_PORT", "6060")
	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Handler:      http.NewCrossOriginProtection().Handler(mux),
	}

	log.Printf("server listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
