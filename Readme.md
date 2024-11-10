### Go Reverse Proxy in Docker for Local Development

This project is a Go-based reverse proxy that allows local development with custom subdomains pointing to specific ports on your local machine. Using Docker, this proxy server dynamically routes requests to different ports based on subdomain configurations.

#### Commands
```
docker build -t go-proxy .
docker run -p 80:80 -v $(pwd)/routes.json:/app/routes.json proxy-server
```
