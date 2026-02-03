# 🤖 Meme TG Bot

Telegram bot written in Go that allows you to download videos by links from popular social networks and platforms using yt-dlp.

The bot supports:

- 📥 Downloading videos from popular platforms (YouTube, TikTok, Instagram, Twitter/X, etc.)

- 🍪 Cookie-based authentication to bypass login and CAPTCHA restrictions

- 💾 SQLite database for storing bot data

- ⚡ Fast processing with in-memory caching

-   🛠 Easy setup with Makefile and environment variables


## 📦 Requirements

Before running the bot make sure you have installed:
- Go 1.25
- yt-dlp
- SQLite
- ffmpeg

Check installed versions:
```bash
go version
yt-dlp --version
sqlite3 --version
ffmpeg -version
````

fedora example installing ffmpeg:
```bash
sudo dnf install https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm
sudo dnf install ffmpeg ffmpeg-devel
```

## ⚙️ Environment variables

Create .env file based on env.example:
```env
TG_BOT_TOKEN=              # REQUIRED - Telegram bot token
LOGGER_BOT_TOKEN=          # OPTIONAL (if used for logging bot)
ADMIN_ID=                  # REQUIRED - Telegram admin user ID
ACTIVATION_KEY=            # REQUIRED - password to use the bot
DB_FILE_NAME=./data/bot.db # REQUIRED - SQLite db file (*.db)
COOKIES_PATH=./data/cookies.txt # REQUIRED - path to cookies file
YTDLP_PATH=/usr/bin/yt-dlp # REQUIRED - yt-dlp binary path
FFMPEG_PATH=/usr/bin/ffmpeg # REQUIRED - ffmpeg binary path
IS_DEBUG=false             # OPTIONAL - print logs for debugging
```
## 📁 Project commands
Makefile included.

### Build:
```bash
make build
```
### Run locally:
```bash
make run
```
### Tidy dependencies:
```bash
make tidy
```

## 🍪 Cookies (IMPORTANT)

Bot uses yt-dlp and requires cookies to bypass CAPTCHA and login restrictions on some websites.
Recommended export via browser extensions:

Chrome:
https://chromewebstore.google.com/detail/get-cookiestxt-locally/cclelndahbckbenkjhflpdbgdldlbecc

Firefox:
https://addons.mozilla.org/en-US/firefox/addon/cookies-txt/

Save exported file as:
```
cookies.txt
```
and set:
```env
COOKIES_PATH=./data/cookies.txt
```
