package vapi

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	yellowPrint = color.New(color.FgYellow).SprintFunc()
	redPrint    = color.New(color.FgRed).SprintFunc()
	greenPrint  = color.New(color.FgGreen, color.Bold).SprintFunc()
	bluePrint   = color.New(color.FgBlue, color.Bold).SprintFunc()
)

//CheckLink (msg string) (string, error)
func CheckLink(msg string) (string, error) {

	fmt.Printf("%s %s \n", yellowPrint("CheckLink>"), greenPrint("msg"))
	if (strings.Contains(msg, "@")) != true && (strings.Contains(msg, "t.me")) != true || msg == "" {
		errMsg := "Не указана ссылка на чат. Требуется ссылка формата `t.me/joinchat/RWltlc4VqhWRzezx` или `@link`"
		err1 := errors.New(errMsg)
		return errMsg, err1

	}
	return "", nil
}
func botAwait(token string) {

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("%s %s %s\n", redPrint("error>"), greenPrint("botAwait>"), yellowPrint(err))

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
		message := update.Message
		switch update.Message.Command() {
		case "id":
			returnChatid(C.bot, message)
		case "check":
		//	checkChat(C, message)
		case "link":
			linkChat(C, message)
		}

		fmt.Printf("%s %s %s %s\n", redPrint("message>"), yellowPrint(message.Chat.ID), bluePrint(message.From.UserName+">"), greenPrint(message.Text))

	}
}

func checkBots() {
	y := color.New(color.FgYellow, color.Bold)
	b := color.New(color.FgBlue, color.Bold)

	g := color.New(color.FgGreen, color.Bold)
	y.Printf("checkBots> ")
	b.Println("start")
	var bots []ClientConfig

	DB.Find(&bots)
	y.Printf("bots len> ")
	g.Println(len(bots))

	for i := 0; i < len(bots); i++ {
		b.Println("bot>", bots[i])
		go botAwait(bots[i].BotToken)
	}
}
