# Stage 1: Build Stage
FROM ubuntu:22.04 AS builder

# Prevent interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install dependencies
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    ffmpeg \
    python3 \
    python3-pip \
    build-essential \
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

# Download and install Golang
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz && \
    rm go1.24.0.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"

# Install yt-dlp via pip
RUN pip3 install --upgrade yt-dlp

WORKDIR /app
# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./

# Source code
#- backend
COPY main.go ./
COPY modules/ modules/

#-frontends
COPY internal/web internal/web
RUN go mod tidy && go mod download


# Build
RUN mkdir -p /compiled && CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-X 'GoRecordurbate/modules/db.Version=$(git describe --tags --always --dirty)'" -a -installsuffix cgo -o /compiled/server main.go

# Run
COPY --from=lunanightbyte/gorecord-base:latest /compiled/server ./server
RUN mkdir -p /app/internal/settings
COPY ./internal/settings/* /app/internal/settings
COPY ./internal/settings/db /app/internal/settings/db
# Expose the port the server listens on
EXPOSE 80

# Start the application
CMD ["./server"]
