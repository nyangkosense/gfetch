package main

import (
	"flag"
	"gfetch/render"
)

func init() {
	flag.StringVar(&render.LogoFlag, "logo", "", "Specify the logo to display (e.g., 'arch', 'debian', 'ubuntu')")
	flag.BoolVar(&render.FullFlag, "full", false, "Display full system information")
}

func main() {
	flag.Parse()
	render.DisplayInfo()
}
