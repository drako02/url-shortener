FROM golang:1.23.3-alpine AS build-stage

WORKDIR /app

# Install system dependencies
RUN apk add --no-cache gcc musl-dev pkgconfig librdkafka-dev


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -o /url-shortener


# FROM gcr.io/distroless/base-debian11 AS build-release-stage
# FROM ubuntu:22.04 AS build-release-stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache librdkafka ca-certificates

WORKDIR /

COPY --from=build-stage  /url-shortener /url-shortener

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 8080

ENTRYPOINT ["/url-shortener"]
