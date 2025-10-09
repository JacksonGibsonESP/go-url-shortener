package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/JacksonGibsonESP/go-url-shortener/internal/config"
	"github.com/JacksonGibsonESP/go-url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func URLRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", ShortCreationHandler)
	r.Get("/{short}", URLHandler)
	return r
}

func ShortCreationHandler(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "text/plain" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	body, _ := io.ReadAll(req.Body)
	url := string(body)
	fmt.Printf("URL requested to short: %s\n", url)

	if strings.TrimSpace(url) == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := service.CreateShortURL(url)
	fmt.Printf("URL shortened: %s\n", shortURL)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.Config.TargetAdress + shortURL))
}

func URLHandler(res http.ResponseWriter, req *http.Request) {
	short := "/" + chi.URLParam(req, "short")
	fmt.Printf("Short URL requested: %s\n", short)

	url := service.GetURLByShort(short)
	fmt.Printf("URL found: %s\n", url)

	if strings.TrimSpace(url) == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
