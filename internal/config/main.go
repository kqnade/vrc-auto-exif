package config

import (
	"encoding/json"
	"fmt"
	"os"
	"log"
	"github.com/fsnotify/fsnotify"
)

type Config struct {
	Enable  bool   `json:"Enable"`
	VRChatPhotoDirectory string `json:"VRChatPhotoDirectory"`
	Debug  bool  `json:"debug"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// When error happend, return a default config
	if config.Debug == true {
		log.SetFlags(log.LstdFlags | log.Lshortfile) // Enable detailed logging
		// Overwrite the default values
		config.Enable = true
		config.VRChatPhotoDirectory = "./test/watch"
		log.Println("Debug mode enabled. Using default configuration values.")
	} else if config.Enable == false {
		log.SetFlags(log.LstdFlags) // Disable detailed logging
		log.Println("Debug mode disabled. Using default configuration values.")
	}

	return &config, nil
}

// SaveConfig saves the configuration to a file.
func SaveConfig(filePath string, config *Config) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config file: %w", err)
	}

	log.Println("Configuration saved to", filePath)
	return nil
}

// UpdateConfig updates the configuration in memory and saves it to the file.
func UpdateConfig(filePath string, updateFunc func(*Config)) error {
	config, err := LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	updateFunc(config)

	if err := SaveConfig(filePath, config); err != nil {
		return fmt.Errorf("failed to save updated config: %w", err)
	}

	log.Println("Configuration updated successfully")
	return nil
}
// WatchConfigFile watches the configuration file for changes and reloads it.
func WatchConfigFile(filePath string, onChange func(*Config)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	err = watcher.Add(filePath)
	if err != nil {
		return fmt.Errorf("failed to add file to watcher: %w", err)
	}
	log.Println("Watching configuration file:", filePath)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("Configuration file changed:", event.Name)
				config, err := LoadConfig(filePath)
				if err != nil {
					log.Println("Error reloading config:", err)
					continue
				}
				onChange(config)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Println("Error watching configuration file:", err)
		}
	}
}
