package handler

import (
	"go-proxy/config"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type errorPageData struct {
	StatusCode  int
	Message     string
	Description string
}

func resolveHostName(fullHostName string) string {
	host, _, err := net.SplitHostPort(fullHostName)
	if err != nil {
		host = fullHostName
	}

	return host
}

func renderErrorPage(writer http.ResponseWriter, statusCode int, message string, description string) {
	writer.WriteHeader(statusCode)

	tmpl, err := template.New("errorPage").ParseFiles("templates/errorPage.html")
	if err != nil {
		http.Error(writer, "Failed to load error page template", http.StatusInternalServerError)
		log.Printf("Error loading template: %v", err)
		return
	}

	errorData := errorPageData{
		StatusCode:  statusCode,
		Message:     http.StatusText(statusCode),
		Description: description,
	}

	err = tmpl.ExecuteTemplate(writer, "errorPage.html", errorData)
	if err != nil {
		http.Error(writer, "Failed to render error page", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func ProxyHandler(writer http.ResponseWriter, request *http.Request) {
	host := resolveHostName(request.Host)

	if host == "localhost" || host == "127.0.0.1" {
		config.RenderConfigPage(writer)
		return
	}

	port, exists := config.GetConfig().Routes[host]

	if !exists {
		log.Printf("Not Found: Host '%s' not in configuration", request.Host)
		renderErrorPage(writer, http.StatusNotFound, "Not Found", "The requested host is not configured to any port")
		return
	}

	// targetUrl := "http://localhost:" + port
	targetUrl := "http://host.docker.internal:" + port // since we are running in docker container and our services are running in local

	backend, err := url.Parse(targetUrl)

	if err != nil {
		http.Error(writer, "Failed to load URL", http.StatusInternalServerError)
		log.Printf("Failed to Parse Backend URL %v", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backend)

	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		log.Printf("Proxy Error: %v", err)
		renderErrorPage(writer, http.StatusInternalServerError, "Internal Server Error", "The backend service encountered an error")
	}

	proxy.ServeHTTP(writer, request)
}
