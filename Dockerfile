FROM golang:1.22.5 AS go

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-mystuff ./cmd/web

FROM debian:stable-slim

WORKDIR /app

COPY --from=go /docker-mystuff /bin/docker-mystuff

EXPOSE 8080

CMD ["/bin/docker-mystuff"]
