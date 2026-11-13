package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	jsxtranspiler "github.com/dimaspandu/vull-stack/internal/jsxtranspiler"
	api_stories "github.com/dimaspandu/vull-stack/services/stories/api"
	web_stories "github.com/dimaspandu/vull-stack/services/stories/web"
)

var assetsDir = http.Dir("client/assets")

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func serveAssets(w http.ResponseWriter, r *http.Request) {
	rel := strings.TrimPrefix(path.Clean(r.URL.Path), "/assets/")

	if strings.HasSuffix(rel, ".jsx") {
		f, err := assetsDir.Open(rel)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		src, err := io.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		out, err := jsxtranspiler.Compile(string(src), "elementBuilder")
		if err != nil {
			http.Error(w, "jsx transpile error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(out))
		return
	}

	http.StripPrefix("/assets/", http.FileServer(assetsDir)).ServeHTTP(w, r)
}

func main() {
	_ = godotenv.Load(".env")
	mux := http.NewServeMux()

	mux.Handle("/assets/", http.HandlerFunc(serveAssets))

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
	mux.HandleFunc("GET /api/stories/{slug}", apiH.Show)

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
