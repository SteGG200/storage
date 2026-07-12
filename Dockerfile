FROM golang:alpine3.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/storage ./cmd/storage
 

FROM alpine:3.24 AS production-stage

RUN apk add --no-cache ca-certificates

RUN mkdir -p /storage
VOLUME /storage

WORKDIR /app

COPY --from=build-stage /app/bin/storage .

EXPOSE 8080

ENTRYPOINT ["/app/storage", "/storage"]
