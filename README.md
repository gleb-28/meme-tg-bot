# 🤖 Meme TG Bot

The Meme TG Bot is a Go-based Telegram bot that downloads videos from popular platforms.

## 🧱 Tech stack

- **Go 1.25** – primary language.
- **Telebot v4** – Telegram framework.
- **Looplab FSM** – FSM.
- **GORM + sqlite driver** – persistence layer for chats and forward-mode settings (`gorm.io/gorm`, `gorm.io/driver/sqlite`).
- **SQLite** – lightweight storage for bot data.
- **yt-dlp** – downloads social-media videos (`github.com/lrstanley/go-ytdlp` wrapping the binary).
- **gallery-dl** – grabs Instagram photos/albums when `yt-dlp` only returns static media; it consumes the same cookies file.
- **ffmpeg** – compresses downloaded videos before sending.
- **cleanenv** – loads `.env` file.

The bot supports:

- 📥 Downloading videos from popular platforms (YouTube, TikTok, Instagram, Twitter/X, etc.)
- 📸 Downloading Instagram photos/albums via `gallery-dl` whenever a video stream isn't available
- 📥 Forward mode to chosen chat
- Forwarding audio, video, video notes, docs, pics, stickers, voices, GIFs, albums in forward mode (albums keep the original order)
- 💾 SQLite database for storing bot data
- ⚡ Fast processing with in-memory caching
- 🛠 Easy setup with Makefile and environment variables
- 🧹 Non-authenticated sessions auto-expire after 10 minutes to keep memory clean (tune `NonAuthSessionTTL` in `internal/core/constants`).

## 🔀 Forward mode

Let the bot forward everything you send it to a single destination chat (group/channel).

### Enable or change the target chat
1. Send `/change_mode` to the bot (after activation).
2. Tap **Включить режим пересылки**. If a previous chat exists, choose whether to reuse it or pick a new one.
3. If picking a new chat, forward any message from the destination group/channel. The bot validates that message comes from a group and that the bot is an **admin** there; otherwise it will ask you to promote it first.
4. On success the bot replies that forwarding is enabled and saves the chat for next time.

### Disable forwarding
1. Send `/change_mode`.
2. Tap **Выключить режим пересылки**.

### How forwarding behaves
- Non-link messages are forwarded; link messages are still downloaded normally.
- Text is prefixed with “<name> говорит: …”; media (photos/videos/video-notes/docs/audio/voice/stickers/GIFs) with “<name> присылает”.
- Albums are buffered briefly (~600 ms) so items stay in order; a single caption is applied. Captions over 1024 characters are sent as a separate text message.
- Avoid pointing the destination chat back to the same conversation where you run the bot to prevent loops.


## 📦 Requirements

Before running the bot make sure you have installed:
- Go 1.25
- yt-dlp
- SQLite
- ffmpeg
- gallery-dl (Python 3 CLI for fetching Instagram photos; install via `pip install --user gallery-dl` or your distro package)

Check installed versions:
```bash
go version
yt-dlp --version
sqlite3 --version
ffmpeg -version
gallery-dl -version
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

Both `yt-dlp` and the new `gallery-dl` Instagram image extractor consume the same cookies file, so keep `COOKIES_PATH` pointed at the export you generated.

## 🚀 VPS Deployment

This guide shows how to deploy the bot on a fresh Ubuntu VPS using Docker.
All deployment assets (compose file, helper script, Dockerfile, and env templates) live under `deploy/`. The Dockerfile now installs `gallery-dl`, so Instagram image downloads keep working inside the container.

1. Create prod.env with and other constants:
```env
DB_FILE_NAME=/app/data/bot.db
COOKIES_PATH=/app/cookies.txt
YTDLP_PATH=/usr/local/bin/yt-dlp
FFMPEG_PATH=/usr/bin/ffmpeg
```

2. Create `deploy/prod.env` (if you need to override the defaults above) and `deploy/deploy.env`, then run the deploy helper from the repo root:
```
sudo chmod +x deploy/deploy.sh && make deploy
```
