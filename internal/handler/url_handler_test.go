package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JacksonGibsonESP/go-url-shortener/internal/service"
)

func TestWebhookPOST(t *testing.T) {
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

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			rr := httptest.NewRecorder()

			Webhook(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedBody != "" {
				if body := rr.Body.String(); body != tt.expectedBody {
					t.Errorf("handler returned unexpected body: got %v want %v",
						body, tt.expectedBody)
				}
			}

			if tt.expectedStatus == http.StatusCreated {
				if contentType := rr.Header().Get("Content-Type"); contentType != "text/plain" {
					t.Errorf("handler returned wrong content type: got %v want text/plain",
						contentType)
				}
			}
		})
	}
}

func TestWebhookGET(t *testing.T) {
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
			expectedStatus:   http.StatusBadRequest,
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

			req := httptest.NewRequest(http.MethodGet, tt.shortURL, nil)

			rr := httptest.NewRecorder()

			Webhook(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if location := rr.Header().Get("Location"); location != tt.expectedLocation {
				t.Errorf("handler returned wrong location: got %v want %v",
					location, tt.expectedLocation)
			}
		})
	}
}

func TestWebhookUnsupportedMethods(t *testing.T) {
	methods := []string{
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			rr := httptest.NewRecorder()

			Webhook(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code for %s: got %v want %v",
					method, status, http.StatusBadRequest)
			}
		})
	}
}
