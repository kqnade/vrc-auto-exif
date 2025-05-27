package watcher

import (
	"fmt"
	"log"
	"time"
	"github.com/fsnotify/fsnotify"
	"os"
)

func WatchDirectory(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
		os.Exit(1)
		return
	}
	defer watcher.Close()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal("Error adding directory to watcher:", err)
		return
	}

	fmt.Println("Watching directory:", dir)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("Created file:", event.Name)
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("Modified file:", event.Name)
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Println("Removed file:", event.Name)
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				fmt.Println("Renamed file:", event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error:", err)
		case <-time.After(time.Second):
		}
	}
}