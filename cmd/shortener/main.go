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
var shortToUrl map[string]string = make(map[string]string)

const shortUrlLength = 8

func createShortUrl(url string) string {
	shortUrl, ok := urlToShort[url]
	if ok {
		return shortUrl
	} else {
		shortUrl = "/" + randomString(shortUrlLength)
		urlToShort[url] = shortUrl
		shortToUrl[shortUrl] = url
		return shortUrl
	}
}

func getUrlByShort(short string) string {
	url, ok := shortToUrl[short]
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

		shortUrl := createShortUrl(url)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://localhost:8080" + shortUrl))
	case http.MethodGet:
		shortUrl := req.URL.Path

		fmt.Println("Short URL requested:")
		fmt.Println(shortUrl)

		url := getUrlByShort(shortUrl)

		fmt.Println("URL found:")
		fmt.Println(url)

		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
