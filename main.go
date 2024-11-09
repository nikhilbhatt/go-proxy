package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Route struct {
	Subdomain string `json:"subdomain"`
	Port      string `json:"port"`
}

type Config struct {
	Routes map[string]string `json:"routes"`
}

var config Config

func loadConfig() {
	file, err := os.Open("routes.json")

	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Error decoding the file: %v", err)
	}

	log.Printf("Loaded config: %+v", config)
}

func handler(w http.ResponseWriter, r *http.Request) {
	port, exists := config.Routes[r.Host]

	if !exists {
		log.Printf("Not Found: Host '%s' not in configuration. Current Routes: %+v", r.Host, config.Routes)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// targetUrl := "http://localhost:" + port
	targetUrl := "http://host.docker.internal:" + port // since we are running in docker container and our services are running in local

	backend, err := url.Parse(targetUrl)

	if err != nil {
		http.Error(w, "Failed to load URL", http.StatusInternalServerError)
		log.Printf("Failed to Parse Backend URL %v", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backend)
	proxy.ServeHTTP(w, r)
}

func main() {
	loadConfig()
	http.HandleFunc("/", handler)
	log.Println("Starting Proxy Server")
	log.Fatal(http.ListenAndServe(":80", nil))
}
