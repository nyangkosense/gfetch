package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	logoFlag string
	fullFlag bool
)

func init() {
	flag.StringVar(&logoFlag, "logo", "", "Specify the logo to display (e.g., 'arch', 'debian', 'ubuntu')")
	flag.BoolVar(&fullFlag, "full", false, "Display full system information")
}

func main() {
	flag.Parse()
	displayInfo()
}

func getWindowManagerAndDE() (string, string) {
	display := os.Getenv("DISPLAY")
	if display == "" {
		return "N/A", "N/A"
	}

	wm := getWindowManager()
	de := getDesktopEnvironment()

	return wm, de
}

func getWindowManager() string {
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

func getDesktopEnvironment() string {
	de := os.Getenv("XDG_CURRENT_DESKTOP")
	if de == "" {
		de = os.Getenv("DESKTOP_SESSION")
	}
	if de == "" {
		de = "Unknown"
	}
	return de
}

func displayInfo() {
	hostInfo, _ := host.Info()
	cpuInfo, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()
	wm, de := getWindowManagerAndDE()

	var asciiArt []string
	if logoFlag != "" {
		asciiArt = getSpecificASCIIArt(logoFlag)
	} else {
		asciiArt = getASCIIArt()
	}
	
	for _, line := range asciiArt {
		fmt.Println(Rainbow(line))
	}

	osName := strings.Title(hostInfo.Platform)
	if hostInfo.PlatformFamily != "" && hostInfo.PlatformFamily != hostInfo.Platform {
		osName = strings.Title(hostInfo.PlatformFamily)
	}

	fmt.Println(ColoredLabel("host", hostInfo.Hostname))
	fmt.Println(ColoredLabel("os", fmt.Sprintf("%s %s", osName, hostInfo.PlatformVersion)))
	fmt.Println(ColoredLabel("packages", getPackageCount(osName)))
	fmt.Println(ColoredLabel("kernel", hostInfo.KernelVersion))
	fmt.Println(ColoredLabel("up", formatUptime(hostInfo.Uptime)))
	
	if fullFlag {
		
		fmt.Println(ColoredLabel("shell", os.Getenv("SHELL")))
	    fmt.Println(ColoredLabel("cpu", cpuInfo[0].ModelName))
	    fmt.Println(ColoredLabel("mem", fmt.Sprintf("%s / %s", formatBytes(memInfo.Used), formatBytes(memInfo.Total))))
		fmt.Println(ColoredLabel("wm", wm))
		fmt.Println(ColoredLabel("de", de))
		fmt.Println(ColoredLabel("gpu", getGPUInfo()))
		fmt.Println(ColoredLabel("disk", getDiskInfo()))
		fmt.Println(ColoredLabel("term", os.Getenv("TERM")))
	}

	fmt.Println(DisplayColorBlocks())
}

func getPackageCount(osName string) string {
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

func getGPUInfo() string {
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

func getDiskInfo() string {
	partitions, err := disk.Partitions(false)
	if err != nil || len(partitions) == 0 {
		return "Unknown"
	}

	usage, err := disk.Usage(partitions[0].Mountpoint)
	if err != nil {
		return "Unknown"
	}

	return fmt.Sprintf("%s / %s", formatBytes(usage.Used), formatBytes(usage.Total))
}

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