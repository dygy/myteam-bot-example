package main

import (
	"context"
	"fmt"
	"github.com/mail-ru-im/bot-golang"
	"log"
	"strings"
	"time"
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
	addTSK := botgolang.NewCallbackButton("add task", "task")
	resBrn := botgolang.NewCallbackButton("Result of the week", "res")
	spamButton := botgolang.NewCallbackButton("Spam this chat", "spam"+account)
	testLinkButton := botgolang.NewURLButton("test", "https://mail.ru/")

	message := bot.NewMessage(account)
	keyboard := botgolang.NewKeyboard()
	keyboard.AddRow(
		addTSK,
		resBrn,
		spamButton,
		testLinkButton,
	)
	message.AttachInlineKeyboard(keyboard)
	message.Text = " "
	if err := message.Send(); err != nil {
		log.Printf("failed to send message: %s", err)
	}
	for update := range updates {
		fmt.Println(update.Type, update.Payload)
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
				addTSK := botgolang.NewCallbackButton("add task", "task")
				resBrn := botgolang.NewCallbackButton("Result of the week", "res")
				spamButton := botgolang.NewCallbackButton("Spam this chat", "spam"+update.Payload.Chat.ID)
				testLinkButton := botgolang.NewURLButton("test", "https://mail.ru/")

				message := bot.NewMessage(update.Payload.Chat.ID)
				keyboard := botgolang.NewKeyboard()
				keyboard.AddRow(
					addTSK,
					resBrn,
					spamButton,
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
			if strings.Contains(data.CallbackData, "spam") {
				spam(bot, strings.Replace(data.CallbackData, "spam", "", 1))
			} else if data.CallbackData == "res" {
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

func spam(bot *botgolang.Bot, account string) {
	for i := 0; i < 500; i++ {
		time.Sleep(100)
		message := bot.NewMessage(account)
		message.Text = "Sed ut perspiciatis, unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam eaque ipsa, quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt, explicabo. nemo enim ipsam voluptatem, quia voluptas sit, aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos, qui ratione voluptatem sequi nesciunt, neque porro quisquam est, qui dolorem ipsum, quia dolor sit, amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt, ut labore et dolore magnam aliquam quaerat voluptatem. ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur? quis autem vel eum iure reprehenderit, qui in ea voluptate velit esse, quam nihil molestiae consequatur, vel illum, qui dolorem eum fugiat, quo voluptas nulla pariatur?"
		err := message.Send()
		if err != nil {
			log.Fatal(err)
		}
	}
}
