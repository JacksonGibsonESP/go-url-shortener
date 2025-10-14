package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/JacksonGibsonESP/go-url-shortener/internal/config"
	"github.com/JacksonGibsonESP/go-url-shortener/internal/logging"
	"github.com/JacksonGibsonESP/go-url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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
	logging.Log.Info("Requested to short", zap.String("url", url))

	if strings.TrimSpace(url) == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := service.CreateShortURL(url)
	logging.Log.Info("Shortened", zap.String("shortUrl", shortURL))

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.Config.TargetAdress + shortURL))
}

func URLHandler(res http.ResponseWriter, req *http.Request) {
	short := "/" + chi.URLParam(req, "short")
	logging.Log.Info("Full URL requested by short", zap.String("shortUrl", short))

	url := service.GetURLByShort(short)
	logging.Log.Info("URL found", zap.String("url", url))

	if strings.TrimSpace(url) == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
