package main

import (
	"context"
	"github.com/mail-ru-im/bot-golang"
	"log"
)

func main() {
	bot, err := botgolang.NewBot(
		GoDotEnvVariable("TOKEN"),
		botgolang.BotApiURL("https://api.internal.myteam.mail.ru/bot/v1"),
	)

	if err != nil {
		log.Fatalf("cannot connect to bot: %s", err)
	}

	log.Println(bot.Info)
	var account = GoDotEnvVariable("account")
	var ctx = context.Background()
	updates := bot.GetUpdatesChannel(ctx)
	var fileNow = ""
	for update := range updates {
		// fmt.Println(update.Type, update.Payload)
		switch update.Type {
		case botgolang.NEW_MESSAGE:
			message := update.Payload.Message()
			if fileNow == "true" {
				fileNow = message.Text
				message := bot.NewMessage(account)
				message.Text = "now type a res"
				err := message.Send()
				if err != nil {
					log.Fatal(err)
				}
			} else if fileNow != "" {
				WriteFile(fileNow, message.Text)
				fileNow = ""
				message := bot.NewMessage(account)
				message.Text = "all done"
				err := message.Send()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				var chat = update.Payload.Chat.ID
				addTSK := botgolang.NewCallbackButton("add task", "task")
				resBrn := botgolang.NewCallbackButton("Result of the week", "res")
				testLinkButton := botgolang.NewURLButton("test", "https://mail.ru/")

				message := bot.NewMessage(chat)
				keyboard := botgolang.NewKeyboard()
				keyboard.AddRow(
					addTSK,
					resBrn,
					testLinkButton,
				)
				message.AttachInlineKeyboard(keyboard)
				message.Text = " "
				if err := message.Send(); err != nil {
					log.Printf("failed to send message: %s", err)
				}
			}
		case botgolang.CALLBACK_QUERY:
			data := update.Payload.CallbackQuery()
			if data.CallbackData == "res" {
				err := bot.SendMessage(bot.NewTextMessage(
					account,
					allFilesData(),
				))
				if err != nil {
					log.Fatal(err)
				} else {
					DeleteAllFiles()
				}
			} else {
				fileNow = "true"
				message := bot.NewMessage(account)
				message.Text = "type task name"
				err := message.Send()
				if err != nil {
					log.Fatal(err)
				}
			}
		}

	}
}
