package vapi

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func botAwait(token string) {
	//	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	g := color.New(color.FgGreen, color.Bold).SprintFunc()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("%s %s %s\n", red("error>"), g("botAwait>"), yellow(err))
		//	red.Printf("  %s\n")
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		//log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)
		fmt.Printf("%s %s %s\n", red("message>"), yellow(update.Message.From.UserName+">"), g(update.Message.Text))
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		//bot.Send(msg)
	}
}
func checkBots() {
	y := color.New(color.FgYellow, color.Bold)
	b := color.New(color.FgBlue, color.Bold)

	g := color.New(color.FgGreen, color.Bold)
	y.Printf("checkBots> ")
	b.Println("start")
	var bots []ClientConfig

	DB.Where("").Find(&bots)
	y.Printf("bots len> ")
	g.Println(len(bots))

	for i := 0; i < len(bots); i++ {
		b.Println("bot>", bots[i])
		go botAwait(bots[i].BotToken)
	}
}
