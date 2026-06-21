package main

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func collectSystem() SystemMetrics {
	percents, _ := cpu.Percent(time.Second, false)
	cpuPercent := 0.0
	if len(percents) > 0 {
		cpuPercent = percents[0]
	}

	memStats, _ := mem.VirtualMemory()

	return SystemMetrics{
		CPUPercent: cpuPercent,
		RAMUsedGB:  float64(memStats.Used) / 1e9,
		RAMTotalGB: float64(memStats.Total) / 1e9,
	}
}
