package main

import (
	"net/http"

	"github.com/JacksonGibsonESP/go-url-shortener/internal/handler"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(handler.Webhook))
}
