package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

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
