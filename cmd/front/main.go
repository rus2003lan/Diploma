package main

import (
	"bytes"
	_ "embed"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	frontendPort = "8000"
	backendURL   = "http://localhost:8080"
)

//go:embed index.html
var index []byte

func main() {
	backend, err := url.Parse(backendURL)
	if err != nil {
		log.Fatal("Failed to parse backend URL:", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(backend)

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		_, _ = io.Copy(w, bytes.NewReader(index))
	})

	log.Printf("Frontend server running on http://localhost:%s", frontendPort)
	log.Fatal(http.ListenAndServe(":"+frontendPort, nil))
}
