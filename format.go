package main

import (
	"fmt"
	"strings"
	"time"
)

func formatUptime(uptime uint64) string {
	duration := time.Duration(uptime) * time.Second
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d days", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hours", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d mins", minutes))
	}

	return strings.Join(parts, ", ")
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
