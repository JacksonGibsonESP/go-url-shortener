package service

import (
	"math/rand"
)

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

func CreateShortURL(url string) string {
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

func GetURLByShort(short string) string {
	url, ok := shortToURL[short]
	if ok {
		return url
	} else {
		return ""
	}
}
