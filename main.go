package main

import (
	"flag"
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
