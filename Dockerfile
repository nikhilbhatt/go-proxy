FROM golang:1.23.2-alpine

WORKDIR /app

COPY . .

RUN go build -o proxy-server .

EXPOSE 80

CMD ["./proxy-server"]

