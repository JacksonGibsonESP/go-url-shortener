package main

import (
	"net/http"

	"github.com/JacksonGibsonESP/go-url-shortener/internal/handler"
)

func main() {
	router := handler.URLRouter()
	http.ListenAndServe(":8080", router)
}
