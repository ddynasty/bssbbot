package main

import (
	"log"

	"github.com/Syfaro/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(teltoken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var updchan tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	updchan.Timeout = 60

	// Creating new update channel
	nupch, err := bot.GetUpdatesChan(updchan)
	if err != nil {
		log.Panic(err)
	}

	// Reading from channel
	for update := range nupch {

		ChatID := update.Message.Chat.ID

		if update.Message.IsCommand() {
			if update.Message.Text == "/congr" {
				reply := "Grac, " + update.Message.From.FirstName
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			}
		}
	}
}
