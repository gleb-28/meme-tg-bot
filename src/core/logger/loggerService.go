package logger

import (
	"log"
	"time"

	t "gopkg.in/telebot.v4"
)

type LoggerService struct {
	isDebug   bool
	loggerBot *t.Bot
	adminID   uint32
}

func (loggerService *LoggerService) Debug(message string) {
	if !loggerService.isDebug {
		return
	}
	log.Println("[DEBUG]: " + message)
}

func (loggerService *LoggerService) Info(message string) {
	log.Println("[INFO]: " + message)
}

func (loggerService *LoggerService) Error(message string) {
	err := "[ERROR]: " + message

	log.Println(err)

	if loggerService.loggerBot != nil {
		_, err := loggerService.loggerBot.Send(&t.User{ID: 1309740174}, err)
		if err != nil {
			log.Println("[ERROR]: ERROR WHILE SENDING TO LOGGER BOT: " + err.Error())
		}
	}
}

func GetLoggerService(isDebug bool, loggerBot *t.Bot, adminID uint32) *LoggerService {
	return &LoggerService{isDebug: isDebug, loggerBot: loggerBot, adminID: adminID}
}

func NewLoggerBot(token string, adminID uint32) *t.Bot {
	if token == "" || adminID == 0 {
		return nil
	}

	loggerBotSettings := t.Settings{
		Token:  token,
		Poller: &t.LongPoller{Timeout: 10 * time.Second}, // TODO
	}

	bot, err := t.NewBot(loggerBotSettings)
	if err != nil {
		panic("Error creating logger bot: " + err.Error())
	}

	return bot
}
