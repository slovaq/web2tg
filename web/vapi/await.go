package vapi

import (
	"fmt"
	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/slovaq/web2tg/web/DAL"
	"log"
	"strconv"
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
		message := update.Message
		switch update.Message.Command() {
		case "id":
			{
				msg := tgbotapi.NewMessage(message.Chat.ID, "ID Чата"+strconv.FormatInt(message.Chat.ID, 10))
				bot.Send(msg)
			}
		case "check":
			{
				var link Link
				msg := tgbotapi.NewMessage(message.Chat.ID, "")
				msg.ParseMode = tgbotapi.ModeMarkdown
				DAL.DB.Where("chat_id = ?", message.Chat.ID).First(&link)
				if (link != Link{}) {
					msg.Text = fmt.Sprintf("Чат привязан. Сводка: \n *ChatID*: %d \n UserLink: %s", link.ChatID, link.UserLink)
				} else {
					msg.Text = fmt.Sprintf("Чат *не привязан*")
				}
				bot.Send(msg)

			}
		case "link":
			{
				var link Link
				msg := tgbotapi.NewMessage(message.Chat.ID, "Линка привязана. ")
				msg.ParseMode = tgbotapi.ModeMarkdown
				args := message.CommandArguments()
				admins, err := bot.GetChatAdministrators(message.Chat.ChatConfig())
				if err != nil {
					msg.Text = "Вероятно что-то пошло не так. Проверьте права админа у бота."
				}
				if !userIsAdmin(*message.From, admins) {
					msg.Text = "У тебя нет прав администратора."
				}
				if args == "" {
					msg.Text = "Не указана ссылка на чат. Требуется ссылка формата `t.me/joinchat/RWltlc4VqhWRzezx`"
				}

				DAL.DB.Where("chat_id = ?", message.Chat.ID).First(&link)
				if (Link{} == link) {
					msg.Text = "Чат еще не привязан. Привязываем."
					final, _ := bot.Send(msg)

					link = Link{
						UserLink: args,
						ChatID:   message.Chat.ID,
					}
					DAL.DB.Create(&link)
					bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
						ChannelUsername: message.Chat.UserName,
						ChatID:          message.Chat.ID,
						MessageID:       final.MessageID,
					})
					msg.Text = "Чат теперь привязан! Проверить можно по команде /check"
				}
				bot.Send(msg)

			}
		}

		//log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)
		fmt.Printf("%s %s %s\n", red("message>"), yellow(message.From.UserName+">"), g(message.Text))
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		//bot.Send(msg)
	}
}
func userIsAdmin(member tgbotapi.User, members []tgbotapi.ChatMember) bool {
	for _, admin := range members {
		if member == *admin.User {
			return true
		}
	}
	return false
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
