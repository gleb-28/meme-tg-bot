package main

import (
	"memetgbot/internal"
	"memetgbot/internal/db"
)

func main() {
	db.InitDB()
	bot.InitBot()
}
