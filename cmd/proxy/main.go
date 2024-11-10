package main

import (
	"go-proxy/config"
	"go-proxy/handler"
	"log"
	"net/http"
)

func main() {
	err := config.LoadConfig("routes.json")

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler.ProxyHandler)

	log.Println("Starting Proxy Server")
	log.Fatal(http.ListenAndServe(":80", nil))
}
