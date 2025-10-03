package service

import (
	"strings"
	"testing"
)

func clearMaps() {
	for k := range urlToShort {
		delete(urlToShort, k)
	}
	for k := range shortToURL {
		delete(shortToURL, k)
	}
}

func TestCreateShortURL(t *testing.T) {
	clearMaps()

	t.Run("should create unique short URL for new URL", func(t *testing.T) {
		url := "https://example.com"
		shortURL := CreateShortURL(url)

		if shortURL == "" {
			t.Error("Expected non-empty short URL")
		}

		if !strings.HasPrefix(shortURL, "/") {
			t.Errorf("Expected short URL to start with '/', got: %s", shortURL)
		}

		if len(shortURL) != shortURLLength+1 { // +1 for "/"
			t.Errorf("Expected short URL length %d, got %d", shortURLLength+1, len(shortURL))
		}
	})

	t.Run("should return same short URL for same long URL", func(t *testing.T) {
		url := "https://example.com/page1"
		shortURL1 := CreateShortURL(url)
		shortURL2 := CreateShortURL(url)

		if shortURL1 != shortURL2 {
			t.Errorf("Expected same short URL for same long URL, got %s and %s", shortURL1, shortURL2)
		}
	})

	t.Run("should create different short URLs for different long URLs", func(t *testing.T) {
		url1 := "https://example.com/page1"
		url2 := "https://example.com/page2"

		shortURL1 := CreateShortURL(url1)
		shortURL2 := CreateShortURL(url2)

		if shortURL1 == shortURL2 {
			t.Error("Expected different short URLs for different long URLs")
		}
	})

	t.Run("should handle empty URL", func(t *testing.T) {
		shortURL := CreateShortURL("")

		if shortURL == "" {
			t.Error("Expected short URL even for empty string")
		}

		if !strings.HasPrefix(shortURL, "/") {
			t.Errorf("Expected short URL to start with '/', got: %s", shortURL)
		}
	})
}

func TestGetURLByShort(t *testing.T) {
	clearMaps()

	t.Run("should return original URL for existing short URL", func(t *testing.T) {
		originalURL := "https://example.com/test"
		shortURL := CreateShortURL(originalURL)

		retrievedURL := GetURLByShort(shortURL)

		if retrievedURL != originalURL {
			t.Errorf("Expected %s, got %s", originalURL, retrievedURL)
		}
	})

	t.Run("should return empty string for non-existing short URL", func(t *testing.T) {
		nonExistingShort := "/nonexisting"
		result := GetURLByShort(nonExistingShort)

		if result != "" {
			t.Errorf("Expected empty string for non-existing short URL, got: %s", result)
		}
	})

	t.Run("should return empty string for empty short URL", func(t *testing.T) {
		result := GetURLByShort("")

		if result != "" {
			t.Errorf("Expected empty string for empty short URL, got: %s", result)
		}
	})

	t.Run("should work with multiple URLs", func(t *testing.T) {
		urls := []string{
			"https://example.com/page1",
			"https://example.com/page2",
			"https://example.com/page3",
		}

		shortURLs := make([]string, len(urls))

		// Create short URLs
		for i, url := range urls {
			shortURLs[i] = CreateShortURL(url)
		}

		// Retrieve original URLs
		for i, shortURL := range shortURLs {
			retrievedURL := GetURLByShort(shortURL)
			if retrievedURL != urls[i] {
				t.Errorf("Expected %s, got %s for short URL %s", urls[i], retrievedURL, shortURL)
			}
		}
	})
}

func TestRandomString(t *testing.T) {
	t.Run("should generate string of correct length", func(t *testing.T) {
		length := 10
		result := randomString(length)

		if len(result) != length {
			t.Errorf("Expected length %d, got %d", length, len(result))
		}
	})

	t.Run("should generate different strings", func(t *testing.T) {
		result1 := randomString(8)
		result2 := randomString(8)

		if result1 == result2 {
			t.Error("Expected different random strings")
		}
	})

	t.Run("should contain only allowed characters", func(t *testing.T) {
		const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		result := randomString(100)

		for _, char := range result {
			if !strings.ContainsRune(allowedChars, char) {
				t.Errorf("Invalid character in random string: %c", char)
			}
		}
	})
}
