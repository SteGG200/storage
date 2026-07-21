FROM golang:alpine3.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux make build
 

FROM alpine:3.24 AS production-stage

RUN apk add --no-cache ca-certificates

RUN mkdir -p /storage
VOLUME /storage

RUN useradd -n -s /usr/bin/bash storage && \
	chown -R storage:storage /storage

USER storage

WORKDIR /app

COPY --from=build-stage /app/bin/storage .

EXPOSE 8080

ENTRYPOINT ["/app/storage", "/storage"]
