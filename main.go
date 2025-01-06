package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	Headers         []string `json:"headers"`
	ExcludePatterns []string `json:"excludePatterns"`
}

func loadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	// Define a flag for the configuration file path
	configPath := flag.String("config", "configuration.json", "Path to the configuration JSON file")
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Search in the directory the program is being called from
	searchDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		os.Exit(1)
	}

	// Compile regex patterns for exclusion
	excludeRegexes := make([]*regexp.Regexp, len(config.ExcludePatterns))
	for i, pattern := range config.ExcludePatterns {
		patternRegex := regexp.MustCompile(strings.ReplaceAll(regexp.QuoteMeta(searchDir+string(os.PathSeparator)+pattern), "\\*", ".*"))
		excludeRegexes[i] = patternRegex
	}

	sensitiveContentFound := false

	err = filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded paths
		for _, excludeRegex := range excludeRegexes {
			if excludeRegex.MatchString(path) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Process only files
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNumber := 1
			for scanner.Scan() {
				line := scanner.Text()
				for _, header := range config.Headers {
					if strings.Contains(line, header) {
						fmt.Printf("Sensitive keyword '%s' found in file '%s' on line %d\n", header, path, lineNumber)
						sensitiveContentFound = true
					}
				}
				lineNumber++
			}

			if err := scanner.Err(); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", searchDir, err)
		os.Exit(1)
	}

	// Exit with code 1 if sensitive content has been detected
	if sensitiveContentFound {
		fmt.Println("Sensitive content found!")
		os.Exit(1)
	}

	// No sensitive content, all clear
	fmt.Println("No sensitive content found.")
	os.Exit(0)
}
