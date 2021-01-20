package vapi

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/slovaq/web2tg/web/DAL"
)

var (
	yellowPrint = color.New(color.FgYellow).SprintFunc()
	redPrint    = color.New(color.FgRed).SprintFunc()
	greenPrint  = color.New(color.FgGreen, color.Bold).SprintFunc()
	bluePrint   = color.New(color.FgBlue, color.Bold).SprintFunc()
)

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

	bot.Debug = false

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
					bot.Send(msg)
					break
				}
				if !userIsAdmin(*message.From, admins) {
					msg.Text = "У тебя нет прав администратора."
					bot.Send(msg)
					break
				}
				errMsg, errLink := CheckLink(args)
				if errLink != nil {
					msg.Text = errMsg
					bot.Send(msg)
					break
				}

				var links []Link
				DAL.DB.Where("chat_id = ?", message.Chat.ID).Find(&links)
				fmt.Printf("%s %s link: %s\n", redPrint("botAwait>"), greenPrint("find link>"), yellowPrint(link))
				link = Link{
					UserLink: args,
					ChatID:   message.Chat.ID,
				}
				if len(links) == 0 {
					msg.Text = "Чат еще не привязан. Привязываем."
					final, _ := bot.Send(msg)
					if result := DB.Create(&link); result.Error != nil {
						fmt.Printf("%s  %s  %d is exists\n", redPrint("error>"), args, uint64(message.Chat.ID))
						msg.Text = "Произошла ошибка при привязке. Попробуйте еще раз или обратитесь к системному администратору!"
					} else {
						bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
							ChannelUsername: message.Chat.UserName,
							ChatID:          message.Chat.ID,
							MessageID:       final.MessageID,
						})
						msg.Text = "Чат теперь привязан! Проверить можно по команде /check"
					}
					bot.Send(msg)
					break
				} else {
					DB.Model(&link).Where("chat_id = ?", message.Chat.ID).Updates(link)
					msg.Text = "Ссылка обновлена! Проверить можно по команде /check"
					bot.Send(msg)
					break
				}

			}
		}

		fmt.Printf("%s %s %s %s\n", redPrint("message>"), yellowPrint(message.Chat.ID), bluePrint(message.From.UserName+">"), greenPrint(message.Text))

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
