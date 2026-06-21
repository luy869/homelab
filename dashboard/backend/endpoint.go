package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

type endpointConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func loadEndpoints() []endpointConfig {
	raw := os.Getenv("ENDPOINTS")
	if raw == "" {
		return nil
	}
	var eps []endpointConfig
	if err := json.Unmarshal([]byte(raw), &eps); err != nil {
		return nil
	}
	return eps
}

func collectEndpoints() []EndpointStatus {
	configs := loadEndpoints()
	if len(configs) == 0 {
		return nil
	}

	results := make([]EndpointStatus, len(configs))
	var wg sync.WaitGroup
	httpClient := &http.Client{Timeout: 5 * time.Second}

	for i, cfg := range configs {
		wg.Add(1)
		go func(i int, cfg endpointConfig) {
			defer wg.Done()
			start := time.Now()
			resp, err := httpClient.Get(cfg.URL)
			latency := time.Since(start).Milliseconds()
			ok := err == nil && resp.StatusCode < 500
			if resp != nil {
				resp.Body.Close()
			}
			results[i] = EndpointStatus{
				Name:      cfg.Name,
				URL:       cfg.URL,
				OK:        ok,
				LatencyMs: latency,
			}
		}(i, cfg)
	}
	wg.Wait()
	return results
}
