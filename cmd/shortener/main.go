package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/JacksonGibsonESP/go-url-shortener/internal/config"
	"github.com/JacksonGibsonESP/go-url-shortener/internal/handler"
	"github.com/JacksonGibsonESP/go-url-shortener/internal/logging"
	"go.uber.org/zap"
)

func main() {
	config.Init()
	flag.Parse()

	value, isPresent := os.LookupEnv("SERVER_ADDRESS")
	if isPresent {
		config.Config.CurrentAdress = value
	}

	value, isPresent = os.LookupEnv("BASE_URL")
	if isPresent {
		config.Config.TargetAdress = value
	}

	if err := logging.Initialize(config.Config.LoggingLevel); err != nil {
		panic(err)
	}

	logging.Log.Info("Running server",
		zap.String("server_address", config.Config.CurrentAdress),
		zap.String("target_address", config.Config.CurrentAdress),
		zap.String("logging_level", config.Config.LoggingLevel))

	router := handler.URLRouter()
	http.ListenAndServe(config.Config.CurrentAdress, logging.WithLogging(router, logging.Log))
}
