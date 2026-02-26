package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//go:embed all:build
var staticFiles embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().UTC().Format(time.RFC3339) + `"}`))
	})

	buildFS, err := fs.Sub(staticFiles, "build")
	if err != nil {
		log.Fatalf("static build files not found: %v", err)
	}

	fileServer := http.FileServer(http.FS(buildFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}

		requestPath := strings.TrimPrefix(r.URL.Path, "/")
		if requestPath == "" {
			requestPath = "index.html"
		}

		if _, statErr := fs.Stat(buildFS, requestPath); statErr == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		fallbackReq := r.Clone(r.Context())
		fallbackReq.URL.Path = "/index.html"
		fileServer.ServeHTTP(w, fallbackReq)
	})

	addr := ":" + port
	log.Printf("server listening on http://localhost%s", addr)
	if listenErr := http.ListenAndServe(addr, mux); listenErr != nil {
		log.Fatalf("server failed: %v", listenErr)
	}
}
