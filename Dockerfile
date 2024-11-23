FROM golang:1.23.2-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod ./
RUN go mod download

COPY . .

EXPOSE 80

CMD ["air", "-c", ".air.toml"]
