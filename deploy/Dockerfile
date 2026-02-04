FROM golang:1.25 AS build
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o meme-bot ./cmd/bot

FROM debian:bookworm-slim AS run
RUN apt update && apt install -y \
    ffmpeg \
    python3 \
    curl \
    sqlite3 \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
        -o /usr/local/bin/yt-dlp && chmod +x /usr/local/bin/yt-dlp
WORKDIR /app
RUN mkdir -p /app/data /app/.cache
COPY --from=build /build/meme-bot .
RUN chmod -R 777 /app
ENV HOME=/app
ENV XDG_CACHE_HOME=/app/.cache
EXPOSE 8080
ENTRYPOINT ["./meme-bot"]