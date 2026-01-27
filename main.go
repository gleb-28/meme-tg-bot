package main

import (
	"log"
	c "memetgbot/src/core/config"

	g "github.com/joho/godotenv"
)

func main() {
	err := g.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config, err := c.GetConfig()
	if err != nil {
		panic("Error getting config: " + err.Error())
	}

	log.Println(config)
}
