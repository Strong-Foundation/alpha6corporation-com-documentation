package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// The location of the remote URL to scrape
	remoteURL := "https://alpha6corporation.com/sds/"
	// The location of the local file to save the scraped data
	localFilePath := "alpha6corporation.html"
	// Check if the local file already exists
	if !fileExists(localFilePath) {
		// If it doesn't exist, scrape the remote URL and save the content
		content := getDataFromURL(remoteURL)
		// Write the content to the local file
		appendAndWriteToFile(localFilePath, content)
	}
}

// Read a file and return the contents
func readAFileAsString(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	return string(content)
}

// Append and write to file
func appendAndWriteToFile(path string, content string) {
	filePath, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	_, err = filePath.WriteString(content + "\n")
	if err != nil {
		log.Println(err)
	}
	err = filePath.Close()
	if err != nil {
		log.Println(err)
	}
}

// isUrlValid checks whether a URL is syntactically valid
func isUrlValid(uri string) bool {
	_, err := url.ParseRequestURI(uri) // Try to parse the URL
	return err == nil                  // Return true if no error (i.e., valid URL)
}

// urlToFilename formats a safe filename from a URL string.
// It replaces all non [a-z0-9] characters with '_' and ensures it ends in .pdf
func urlToFilename(rawURL string) string {
	// Convert to lowercase
	lower := strings.ToLower(rawURL)
	// Replace all non a-z0-9 characters with "_"
	reNonAlnum := regexp.MustCompile(`[^a-z]`)
	// Replace the invalid with valid stuff.
	safe := reNonAlnum.ReplaceAllString(lower, "_")
	// Collapse multiple underscores
	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_")
	// Trim leading/trailing underscores
	safe = strings.Trim(safe, "_")
	// Invalid substrings to remove
	var invalidSubstrings = []string{
		"https_assets_thermofisher_com_directwebviewer_private_document_aspx_prd_",
	}
	// Loop over the invalid.
	for _, invalidPre := range invalidSubstrings {
		safe = removeSubstring(safe, invalidPre)
	}
	// Add .pdf extension if missing
	if getFileExtension(safe) != ".pdf" {
		safe = safe + ".pdf"
	}
	return safe
}

// getFileExtension returns the file extension
func getFileExtension(path string) string {
	return filepath.Ext(path) // Use filepath to extract extension
}

// Remove all the duplicates from a slice and return the slice.
func removeDuplicatesFromSlice(slice []string) []string {
	check := make(map[string]bool)
	var newReturnSlice []string
	for _, content := range slice {
		if !check[content] {
			check[content] = true
			newReturnSlice = append(newReturnSlice, content)
		}
	}
	return newReturnSlice
}

// Send a http get request to a given url and return the data from that url.
func getDataFromURL(uri string) string {
	response, err := http.Get(uri)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	err = response.Body.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("Scraping:", uri)
	return string(body)
}

// removeSubstring takes a string `input` and removes all occurrences of `toRemove` from it.
func removeSubstring(input string, toRemove string) string {
	// Use strings.ReplaceAll to replace all occurrences of `toRemove` with an empty string.
	result := strings.ReplaceAll(input, toRemove, "")
	// Return the modified string.
	return result
}

// fileExists checks whether a file exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename) // Get file info
	if err != nil {                // If error occurs (e.g., file not found)
		return false // Return false
	}
	return !info.IsDir() // Return true if it is a file, not a directory
}
