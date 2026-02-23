package tts

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	googleTTSURL = "https://translate.google.com/translate_tts"
	maxTextLen   = 200 // Google TTS limit per request
)

// Generate creates an MP3 file from the given text using Google TTS
// Returns the path to the generated audio file
func Generate(text string) (string, error) {
	// Create temp file
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "cast-tts-*.mp3")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Build URL
	params := url.Values{}
	params.Set("ie", "UTF-8")
	params.Set("client", "tw-ob")
	params.Set("tl", "en")
	params.Set("q", text)

	requestURL := fmt.Sprintf("%s?%s", googleTTSURL, params.Encode())

	// Make request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to look like a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://translate.google.com/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", fmt.Errorf("TTS request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", fmt.Errorf("TTS request returned status %d", resp.StatusCode)
	}

	// Copy audio data to temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to write audio: %w", err)
	}

	tmpFile.Close()
	return tmpPath, nil
}

// Cleanup removes the temporary audio file
func Cleanup(path string) {
	if path != "" && filepath.Dir(path) == os.TempDir() {
		os.Remove(path)
	}
}
