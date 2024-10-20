package sys

import (
	"fmt"
	"gfetch/format"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
)

func GetWindowManagerAndDE() (string, string) {
	display := os.Getenv("DISPLAY")
	if display == "" {
		return "N/A", "N/A"
	}

	wm := GetWindowManager()
	de := GetDesktopEnvironment()

	return wm, de
}

func GetWindowManager() string {
	// Try to get window manager using xprop
	cmd := exec.Command("xprop", "-root", "-notype", "_NET_SUPPORTING_WM_CHECK")
	output, err := cmd.Output()
	if err == nil {
		parts := strings.Split(string(output), " ")
		if len(parts) > 0 {
			id := strings.TrimSpace(parts[len(parts)-1])
			cmd = exec.Command("xprop", "-id", id, "-notype", "-len", "25", "-f", "_NET_WM_NAME", "8t")
			output, err = cmd.Output()
			if err == nil {
				if strings.Contains(string(output), "_NET_WM_NAME = ") {
					wm := strings.Split(string(output), "\"")
					if len(wm) > 1 {
						return wm[1]
					}
				}
			}
		}
	}

	// Fallback to process list checking
	cmd = exec.Command("ps", "x")
	output, err = cmd.Output()
	if err == nil {
		processes := strings.Split(string(output), "\n")
		for _, process := range processes {
			switch {
			case strings.Contains(process, "catwm"):
				return "catwm"
			case strings.Contains(process, "fvwm"):
				return "fvwm"
			case strings.Contains(process, "dwm"):
				return "dwm"
			case strings.Contains(process, "2bwm"):
				return "2bwm"
			case strings.Contains(process, "monsterwm"):
				return "monsterwm"
			case strings.Contains(process, "wmaker"):
				return "Window Maker"
			case strings.Contains(process, "sowm"):
				return "sowm"
			case strings.Contains(process, "penrose"):
				return "penrose"
			}
		}
	}

	return "Unknown"
}

func GetDesktopEnvironment() string {
	de := os.Getenv("XDG_CURRENT_DESKTOP")
	if de == "" {
		de = os.Getenv("DESKTOP_SESSION")
	}
	if de == "" {
		de = "Unknown"
	}
	return de
}

func GetPackageCount(osName string) string {
	var cmd *exec.Cmd
	switch strings.ToLower(osName) {
	case "ubuntu", "debian", "linux mint":
		cmd = exec.Command("dpkg-query", "-f", "${binary:Package}\n", "-W")
	case "fedora", "centos":
		cmd = exec.Command("rpm", "-qa")
	case "arch", "manjaro":
		cmd = exec.Command("pacman", "-Q")
	case "void":
		cmd = exec.Command("xbps-query", "-l")
	default:
		return "Unknown"
	}

	output, err := cmd.Output()
	if err != nil {
		return "Error"
	}

	count := len(strings.Split(strings.TrimSpace(string(output)), "\n"))
	return fmt.Sprintf("%d", count)
}

func GetGPUInfo() string {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("lspci", "-v")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "VGA") || strings.Contains(line, "3D") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						return strings.TrimSpace(parts[1])
					}
				}
			}
		}
	}
	return "Unknown"
}

func GetDiskInfo() string {
	partitions, err := disk.Partitions(false)
	if err != nil || len(partitions) == 0 {
		return "Unknown"
	}

	usage, err := disk.Usage(partitions[0].Mountpoint)
	if err != nil {
		return "Unknown"
	}

	return fmt.Sprintf("%s / %s", format.FormatBytes(usage.Used), format.FormatBytes(usage.Total))
}
