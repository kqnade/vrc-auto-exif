package main

import (
	"fmt"
	"k4na.de/vrc-auto-exif/internal/watcher"
	"k4na.de/vrc-auto-exif/internal/config"
)

func main() {
	fmt.Println("VRC Auto EXIF Editor")

	// Load configuration
	cfg, err := config.LoadConfig("./config.json")
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}
	fmt.Println("Configuration loaded successfully")
	fmt.Printf("Watching directory: %s\n", cfg.VRChatPhotoDirectory)
	// Initialize the watcher
	directory := cfg.VRChatPhotoDirectory
	watcher.WatchDirectory(directory)
	fmt.Println("Watching for new files...")
}
