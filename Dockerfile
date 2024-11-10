FROM golang:1.23.2-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

EXPOSE 80

CMD ["go", "run", "."]
