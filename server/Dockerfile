FROM golang:1.23.1-alpine3.19

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 8080

CMD ["go", "run", "main.go"]