package handler

import (
	"encoding/json"
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
	//RESTfull
	r.Post("/api/shorten", RESTShortCreationHandler)
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

type RequestBody struct {
	Url string `json:"url"`
}

type ResponseBody struct {
	Result string `json:"result"`
}

func RESTShortCreationHandler(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var request RequestBody
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&request); err != nil {
		logging.Log.Error("Error encoding request", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	logging.Log.Info("Requested to short", zap.String("url", request.Url))

	if strings.TrimSpace(request.Url) == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := service.CreateShortURL(request.Url)
	logging.Log.Info("Shortened", zap.String("shortUrl", shortURL))

	response := ResponseBody{
		Result: config.Config.TargetAdress + shortURL,
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(res)
	if err := enc.Encode(response); err != nil {
		logging.Log.Error("Error encoding response", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}
