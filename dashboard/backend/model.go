package main

import "time"

type Status struct {
	Containers []ContainerStatus `json:"containers"`
	Endpoints  []EndpointStatus  `json:"endpoints"`
	System     SystemMetrics     `json:"system"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type ContainerStatus struct {
	Name   string `json:"name"`
	State  string `json:"state"`
	Status string `json:"status"`
}

type EndpointStatus struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	OK        bool   `json:"ok"`
	LatencyMs int64  `json:"latency_ms"`
}

type SystemMetrics struct {
	CPUPercent float64 `json:"cpu_percent"`
	RAMUsedGB  float64 `json:"ram_used_gb"`
	RAMTotalGB float64 `json:"ram_total_gb"`
}
