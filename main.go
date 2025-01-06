package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// TODO: read those from a configuration file, probably in JSON
	// or YAML if I'm feeling brave.
	headers := []string{
		"AGE-SECRET-KEY",              // agenix secret keys
		"BEGIN OPENSSH PRIVATE KEY",   // I don't think I can possibly commit those by accident but still
		"BEGIN PGP PRIVATE KEY BLOCK", // "
		"PRIVATE",                     // Not sure if I have a file with "PRIVATE" as a header, but no harm including this
	}

	excludePatterns := []string{
		"modules/roles/server/system/services/forgejo.nix", // PRIVATE is an infix here
		".git", // avoid parsing the git history
		"*.go", // avoid false positives in `go run`
	}

	// Search in the directory the program is being called from
	searchDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		os.Exit(1)
	}

	// Compile regex patterns for exclusion
	excludeRegexes := make([]*regexp.Regexp, len(excludePatterns))
	for i, pattern := range excludePatterns {
		patternRegex := regexp.MustCompile(strings.ReplaceAll(regexp.QuoteMeta(searchDir+pattern), "\\*", ".*"))
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
				for _, header := range headers {
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
