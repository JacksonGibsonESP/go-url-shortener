package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/JacksonGibsonESP/go-url-shortener/internal/config"
	"github.com/JacksonGibsonESP/go-url-shortener/internal/handler"
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

	router := handler.URLRouter()
	http.ListenAndServe(config.Config.CurrentAdress, router)
}
