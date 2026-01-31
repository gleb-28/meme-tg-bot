# Dockerfile for a Go Telegram Bot

# Use the official Golang image as a base image
FROM golang:1.25-alpine as build

# Install build dependencies (включая git)
RUN apk update && apk add --no-cache git python3 py3-pip ffmpeg

# Set the current working directory inside the container
WORKDIR /build

# Copy go mod and sum files
COPY go.mod .
COPY go.sum .
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source code into the container
COPY . .
# Build the Go application
RUN go build -o main .


# Use a minimal Alpine image for the final stage to create a smaller image
FROM alpine:3.19 as run

RUN apk add --no-cache yt-dlp

# Set the current working directory inside the container
WORKDIR /run

# Copy the compiled application from the builder stage
COPY --from=build /build/main /run/

ENV YT_DLP_PATH="/usr/bin/yt-dlp"

# Expose port 8080 for the application
EXPOSE 8080

# Run the application
ENTRYPOINT ["/run/main"]
