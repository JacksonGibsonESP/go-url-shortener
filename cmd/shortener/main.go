package main

import (
	"flag"
	"net/http"

	"github.com/JacksonGibsonESP/go-url-shortener/internal/config"
	"github.com/JacksonGibsonESP/go-url-shortener/internal/handler"
)

func main() {
	config.Init()
	flag.Parse()

	router := handler.URLRouter()
	http.ListenAndServe(config.Config.CurrentAdress, router)
}
