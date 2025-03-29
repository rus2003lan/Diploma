package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	frontendPort = "8000"                  // Порт для фронтенда (изменили на 8000)
	backendURL   = "http://localhost:8080" // Адрес вашего бекенда
)

func main() {
	// Настройка прокси к бекенду
	backend, err := url.Parse(backendURL)
	if err != nil {
		log.Fatal("Failed to parse backend URL:", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(backend)

	// Модифицируем директиву для API
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// Добавляем CORS заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Проксируем запрос к бекенду
		proxy.ServeHTTP(w, r)
	})

	// Обработчик для фронтенда
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		file, err := os.Open("./index.html")
		if err != nil {
			http.Error(w, fmt.Errorf("Frontend not found: %w", err).Error(), http.StatusNotFound)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "text/html")
		_, _ = io.Copy(w, file)
	})

	log.Printf("Frontend server running on http://localhost:%s", frontendPort)
	log.Fatal(http.ListenAndServe(":"+frontendPort, nil))
}
