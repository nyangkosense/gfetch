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

	osName := strings.Title(hostInfo.Platform)
	if hostInfo.PlatformFamily != "" && hostInfo.PlatformFamily != hostInfo.Platform {
		osName = strings.Title(hostInfo.PlatformFamily)
    	}

	infoLines := []string{
	color.ColoredLabel("host", hostInfo.Hostname),
	color.ColoredLabel("os", fmt.Sprintf("%s %s", osName, hostInfo.PlatformVersion)),
	color.ColoredLabel("wm", wm),
	color.ColoredLabel("packages", sys.GetPackageCount(osName)),
        color.ColoredLabel("kernel", hostInfo.KernelVersion),
        color.ColoredLabel("up", format.FormatUptime(hostInfo.Uptime)),
	}

	if FullFlag {
		additionalInfo := []string{
		color.ColoredLabel("shell", os.Getenv("SHELL")),
		color.ColoredLabel("cpu", cpuInfo[0].ModelName),
		color.ColoredLabel("mem", fmt.Sprintf("%s / %s", format.FormatBytes(memInfo.Used), format.FormatBytes(memInfo.Total))),
		color.ColoredLabel("de", de),
		color.ColoredLabel("gpu", sys.GetGPUInfo()),
		color.ColoredLabel("disk", sys.GetDiskInfo()),
		color.ColoredLabel("term", os.Getenv("TERM")),
		}
	infoLines = append(infoLines, additionalInfo...)
	}

	maxArtWidth := 0
	for _, line := range asciiArt {
		if len(color.StripANSI(line)) > maxArtWidth {
			maxArtWidth = len(color.StripANSI(line))
		}
	}

	maxLines := len(asciiArt)
		if len(infoLines) > maxLines {
        		maxLines = len(infoLines)
		}

    for len(asciiArt) < maxLines {
        asciiArt = append(asciiArt, "")
    }
    for len(infoLines) < maxLines {
        infoLines = append(infoLines, "")
    }

	padding := 4 
	for i := 0; i < maxLines; i++ {
		artLine := ""

	if i < len(asciiArt) {
		artLine = color.Rainbow(asciiArt[i])
	}
        
	infoLine := ""
	if i < len(infoLines) {
		infoLine = infoLines[i]
	}
	paddedArtLine := fmt.Sprintf("%-*s", maxArtWidth+len(artLine)-len(color.StripANSI(artLine)), artLine)
        
	fmt.Printf("%s%s%s\n", paddedArtLine, strings.Repeat(" ", padding), infoLine)
	}

	if len(asciiArt) > 0 {
		fmt.Printf("%s%s%s\n", 
		strings.Repeat(" ", maxArtWidth), 
		strings.Repeat(" ", padding), 
		color.DisplayColorBlocks())
    	} else {
	fmt.Println(color.DisplayColorBlocks())
    }
}
