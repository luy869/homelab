package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

//go:embed all:static
var staticFiles embed.FS

type server struct {
	mu     sync.RWMutex
	status *Status
}

func main() {
	project := os.Getenv("COMPOSE_PROJECT")
	if project == "" {
		project = "homelab"
	}

	srv := &server{}
	srv.refresh(project)

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			srv.refresh(project)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", srv.handleStatus)

	static, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(static)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func (s *server) refresh(project string) {
	var (
		containers []ContainerStatus
		endpoints  []EndpointStatus
		system     SystemMetrics
		wg         sync.WaitGroup
	)

	wg.Add(3)
	go func() { defer wg.Done(); containers = collectContainers(project) }()
	go func() { defer wg.Done(); endpoints = collectEndpoints() }()
	go func() { defer wg.Done(); system = collectSystem() }()
	wg.Wait()

	s.mu.Lock()
	s.status = &Status{
		Containers: containers,
		Endpoints:  endpoints,
		System:     system,
		UpdatedAt:  time.Now(),
	}
	s.mu.Unlock()
}

func (s *server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	status := s.status
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(status)
}
