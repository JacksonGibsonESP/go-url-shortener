package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(webhook))
}

func randomString(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var urlToShort map[string]string = make(map[string]string)
var shortToURL map[string]string = make(map[string]string)

const shortURLLength = 8

func createShortURL(url string) string {
	shortURL, ok := urlToShort[url]
	if ok {
		return shortURL
	} else {
		shortURL = "/" + randomString(shortURLLength)
		urlToShort[url] = shortURL
		shortToURL[shortURL] = url
		return shortURL
	}
}

func getURLByShort(short string) string {
	url, ok := shortToURL[short]
	if ok {
		return url
	} else {
		return ""
	}
}

func webhook(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		if req.Header.Get("Content-Type") != "text/plain" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		body, _ := io.ReadAll(req.Body)
		url := string(body)
		fmt.Println("URL requested to short:")
		fmt.Println(url)

		shortURL := createShortURL(url)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://localhost:8080" + shortURL))
	case http.MethodGet:
		shortURL := req.URL.Path

		fmt.Println("Short URL requested:")
		fmt.Println(shortURL)

		url := getURLByShort(shortURL)

		fmt.Println("URL found:")
		fmt.Println(url)

		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
