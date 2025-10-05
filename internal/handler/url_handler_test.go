package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JacksonGibsonESP/go-url-shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string,
	body string, bodyContentType string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	if bodyContentType != "" {
		req.Header.Set("Content-Type", bodyContentType)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestWebhookPOST(t *testing.T) {
	server := httptest.NewServer(URLRouter())
	defer server.Close()

	tests := []struct {
		name           string
		contentType    string
		body           string
		mockShortURL   string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful URL shortening",
			contentType:    "text/plain",
			body:           "https://example.com",
			mockShortURL:   "/abc123",
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://localhost:8080/abc123",
		},
		{
			name:           "empty URL",
			contentType:    "text/plain",
			body:           "",
			mockShortURL:   "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "wrong content type",
			contentType:    "application/json",
			body:           `{"url": "https://example.com"}`,
			mockShortURL:   "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "missing content type",
			contentType:    "",
			body:           "https://example.com",
			mockShortURL:   "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalCreateShortURLFunction := service.CreateShortURL
			originalGetURLByShort := service.GetURLByShort

			defer func() {
				service.CreateShortURL = originalCreateShortURLFunction
				service.GetURLByShort = originalGetURLByShort
			}()

			service.CreateShortURL = func(url string) string {
				if url != tt.body {
					t.Errorf("Expected URL %s, got %s", tt.body, url)
				}
				return tt.mockShortURL
			}

			resp, rbd := testRequest(t, server, http.MethodPost, "/", tt.body, tt.contentType)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if status := resp.StatusCode; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedBody != "" {
				if rbd != tt.expectedBody {
					t.Errorf("handler returned unexpected body: got %v want %v",
						rbd, tt.expectedBody)
				}
			}

			if tt.expectedStatus == http.StatusCreated {
				if contentType := resp.Header.Get("Content-Type"); contentType != "text/plain" {
					t.Errorf("handler returned wrong content type: got %v want text/plain",
						contentType)
				}
			}
		})
	}
}

func TestWebhookGET(t *testing.T) {
	server := httptest.NewServer(URLRouter())
	defer server.Close()

	tests := []struct {
		name             string
		shortURL         string
		mockOriginalURL  string
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "successful redirect",
			shortURL:         "/abc123",
			mockOriginalURL:  "https://example.com",
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "https://example.com",
		},
		{
			name:             "empty short URL",
			shortURL:         "/",
			mockOriginalURL:  "",
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedLocation: "",
		},
		{
			name:             "non-existent short URL",
			shortURL:         "/nonexistent",
			mockOriginalURL:  "",
			expectedStatus:   http.StatusBadRequest,
			expectedLocation: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			originalCreateShortURLFunction := service.CreateShortURL
			originalGetURLByShort := service.GetURLByShort

			defer func() {
				service.CreateShortURL = originalCreateShortURLFunction
				service.GetURLByShort = originalGetURLByShort
			}()

			service.GetURLByShort = func(short string) string {
				if short != tt.shortURL {
					t.Errorf("Expected short URL %s, got %s", tt.shortURL, short)
				}
				return tt.mockOriginalURL
			}

			resp, _ := testRequest(t, server, http.MethodGet, tt.shortURL, "", "")
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if status := resp.StatusCode; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if location := resp.Header.Get("Location"); location != tt.expectedLocation {
				t.Errorf("handler returned wrong location: got %v want %v",
					location, tt.expectedLocation)
			}
		})
	}
}
