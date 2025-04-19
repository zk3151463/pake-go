package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zk3151463/pake-go/pkg/builder"
	"github.com/zk3151463/pake-go/pkg/config"
)

func main() {
	// Define command line flags
	url := flag.String("url", "", "URL to package")
	name := flag.String("name", "", "Application name")
	icon := flag.String("icon", "", "Application icon path")
	width := flag.Int("width", 1200, "Window width")
	height := flag.Int("height", 800, "Window height")
	hideTitleBar := flag.Bool("hide-title-bar", false, "Hide title bar")
	transparent := flag.Bool("transparent", false, "Enable transparent window")
	alwaysOnTop := flag.Bool("always-on-top", false, "Keep window always on top")
	userAgent := flag.String("user-agent", "", "Custom user agent")
	configFile := flag.String("config", "", "Path to config file")

	flag.Parse()

	// If no URL is provided, show usage
	if *url == "" && *configFile == "" {
		fmt.Println("Usage: pake-go -url <url> [options]")
		fmt.Println("   or: pake-go -config <config-file>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create config
	cfg := &config.Config{
		URL:          *url,
		Name:         *name,
		Icon:         *icon,
		Width:        *width,
		Height:       *height,
		HideTitleBar: *hideTitleBar,
		Transparent:  *transparent,
		AlwaysOnTop:  *alwaysOnTop,
		UserAgent:    *userAgent,
	}

	// If config file is provided, load it
	if *configFile != "" {
		loadedConfig, err := config.LoadConfig(*configFile)
		if err != nil {
			fmt.Printf("Error loading config file: %v\n", err)
			os.Exit(1)
		}
		cfg = loadedConfig
	}

	// Create builder
	b := builder.NewBuilder(cfg)

	// Build the application
	if err := b.Build(); err != nil {
		fmt.Printf("Error building application: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Application built successfully!")
}
