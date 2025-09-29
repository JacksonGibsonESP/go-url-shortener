package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/JacksonGibsonESP/go-url-shortener/internal/service"
)

func Webhook(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		if req.Header.Get("Content-Type") != "text/plain" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		body, _ := io.ReadAll(req.Body)
		url := string(body)
		fmt.Printf("URL requested to short: %s\n", url)

		shortURL := service.CreateShortURL(url)
		fmt.Printf("URL shortened: %s\n", shortURL)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://localhost:8080" + shortURL))
	case http.MethodGet:
		shortURL := req.URL.Path
		fmt.Printf("Short URL requested: %s\n", shortURL)

		url := service.GetURLByShort(shortURL)
		fmt.Printf("URL found: %s\n", url)

		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
