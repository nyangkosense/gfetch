package render

import (
	"fmt"
	"gfetch/art"
	"gfetch/color"
	"gfetch/format"
	"gfetch/sys"
	"os"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	LogoFlag string
	FullFlag bool
)

func DisplayInfo() {
	hostInfo, _ := host.Info()
	cpuInfo, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()
	wm, de := sys.GetWindowManagerAndDE()

	var asciiArt []string
	if LogoFlag != "" {
		asciiArt = art.GetSpecificASCIIArt(LogoFlag)
	} else {
		asciiArt = art.GetASCIIArt()
	}

	for _, line := range asciiArt {
		fmt.Println(color.Rainbow(line))
	}

	osName := strings.Title(hostInfo.Platform)
	if hostInfo.PlatformFamily != "" && hostInfo.PlatformFamily != hostInfo.Platform {
		osName = strings.Title(hostInfo.PlatformFamily)
	}

	fmt.Println(color.ColoredLabel("host", hostInfo.Hostname))
	fmt.Println(color.ColoredLabel("os", fmt.Sprintf("%s %s", osName, hostInfo.PlatformVersion)))
	fmt.Println(color.ColoredLabel("packages", sys.GetPackageCount(osName)))
	fmt.Println(color.ColoredLabel("kernel", hostInfo.KernelVersion))
	fmt.Println(color.ColoredLabel("up", format.FormatUptime(hostInfo.Uptime)))

	if FullFlag {

		fmt.Println(color.ColoredLabel("shell", os.Getenv("SHELL")))
		fmt.Println(color.ColoredLabel("cpu", cpuInfo[0].ModelName))
		fmt.Println(color.ColoredLabel("mem", fmt.Sprintf("%s / %s", format.FormatBytes(memInfo.Used), format.FormatBytes(memInfo.Total))))
		fmt.Println(color.ColoredLabel("wm", wm))
		fmt.Println(color.ColoredLabel("de", de))
		fmt.Println(color.ColoredLabel("gpu", sys.GetGPUInfo()))
		fmt.Println(color.ColoredLabel("disk", sys.GetDiskInfo()))
		fmt.Println(color.ColoredLabel("term", os.Getenv("TERM")))
	}

	fmt.Println(color.DisplayColorBlocks())
}
