package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	TgBotToken     string `env:"TG_BOT_TOKEN" env-required:"true"`
	LoggerBotToken string `env:"LOGGER_BOT_TOKEN"`
	AdminID        int64  `env:"ADMIN_ID" env-required:"true"`
	ActivationKey  string `env:"ACTIVATION_KEY" env-required:"true"`
	CookiesPath    string `env:"COOKIES_PATH" env-required:"true"`
	YtdlpPath      string `env:"YTDLP_PATH" env-required:"true"`
	FfmpegPath     string `env:"FFMPEG_PATH" env-required:"true"`
	GalleryDLPath  string `env:"GALLERYDL_PATH" env-default:"gallery-dl"`
	Instagram      InstagramConfig
	IsDebug        bool           `env:"IS_DEBUG" env-default:"false"`
	Database       DatabaseConfig `env-required:"true"`
}

type InstagramConfig struct {
	GalleryDLUserAgent      string `env:"INSTAGRAM_GDL_USER_AGENT" env-default:"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"`
	GalleryDLSleep          string `env:"INSTAGRAM_GDL_SLEEP" env-default:"7.0-12.0"`
	GalleryDLSleepRequest   string `env:"INSTAGRAM_GDL_SLEEP_REQUEST" env-default:"25.0-50.0"`
	GalleryDLSleepExtractor string `env:"INSTAGRAM_GDL_SLEEP_EXTRACTOR" env-default:"4.0-8.0"`
	GalleryDLSleep429       string `env:"INSTAGRAM_GDL_SLEEP_429" env-default:"900.0"`
}

type DatabaseConfig struct {
	FileName string `env:"DB_FILE_NAME" env-required:"true"`
}

func MustConfig() *AppConfig {
	config := &AppConfig{}

	if os.Getenv("IS_DOCKERIZED") != "true" {
		if err := cleanenv.ReadConfig(".env", config); err != nil {
			log.Fatal("Error loading .env file (might be missing in non-dockerized env): ", err.Error())
		}
	}

	if err := cleanenv.ReadEnv(config); err != nil {
		log.Fatal("Error reading environment variables: ", err.Error())
	}

	return config
}
