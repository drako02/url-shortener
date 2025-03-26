FROM golang:1.23.3 AS build-stage

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    git \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /url-shortener

# EXPOSE 8080

# CMD [ "/url-shortener" ]

# FROM gcr.io/distroless/base-debian11 AS build-release-stage
FROM ubuntu:22.04 AS build-release-stage

WORKDIR /

COPY --from=build-stage  /url-shortener /url-shortener

# Copy shared libraries needed for runtime
COPY --from=build-stage /usr/lib/*/librdkafka*.so* /usr/lib/

# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*


EXPOSE 8080
# Create nonroot user and group
RUN groupadd -r nonroot && useradd -r -g nonroot nonroot

USER nonroot:nonroot

ENTRYPOINT ["/url-shortener"]
