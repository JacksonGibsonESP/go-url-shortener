package config

import (
	"flag"
)

type Params struct {
	CurrentAdress string
	TargetAdress  string
	LoggingLevel  string
}

var Config Params

func Init() {
	flag.StringVar(&Config.CurrentAdress, "a", "localhost:8080", "HTTP Server adress")
	flag.StringVar(&Config.TargetAdress, "b", "http://localhost:8080", "Shortened URL server adress")
	flag.StringVar(&Config.LoggingLevel, "l", "info", "logging level")
}
