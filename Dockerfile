# Build stage
FROM golang:alpine3.24 AS build-stage

RUN apk add --no-cache make git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux make build

# Final runtime stage
FROM alpine:3.24 AS production-stage

RUN apk add --no-cache ca-certificates tzdata

RUN mkdir -p /storage

RUN addgroup -S storage && \
    adduser -S -G storage storage && \
    chown -R storage:storage /storage

USER storage

WORKDIR /app

COPY --from=build-stage /app/bin/storage .

EXPOSE 8080

VOLUME ["/storage"]

ENTRYPOINT ["/app/storage", "/storage"]
