package main

import (
	"net/http"

	"github.com/ZeroTechh/AuthServer/handler"
	"github.com/ZeroTechh/VelocityCore/logger"
	"github.com/ZeroTechh/hades"
)

var (
	config = hades.GetConfig(
		"main.yaml",
		[]string{"config", "../config", "../../config"},
	)
	log = logger.GetLogger(
		config.Map("service").Str("logFile"),
		config.Map("service").Bool("debug"),
	)
)

func main() {
	handler := handler.New()
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/auth", handler.Auth)
	http.HandleFunc("/verify", handler.Verify)
	err := http.ListenAndServe(config.Map("service").Str("address"), nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
