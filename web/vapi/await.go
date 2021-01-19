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

func CheckLink(msg string) (string, error) {
	yellow := color.New(color.FgYellow).SprintFunc()
	//red := color.New(color.FgRed).SprintFunc()
	g := color.New(color.FgGreen, color.Bold).SprintFunc()

	//	fmt.Printf("%s %s %s\n", red("error>"), g("botAwait>"), yellow(err))
	fmt.Printf("%s %s \n", yellow("CheckLink>"), g("msg"))
	if (strings.Contains(msg, "@")) != true && (strings.Contains(msg, "t.me")) != true || msg == "" {
		errMsg := "Не указана ссылка на чат. Требуется ссылка формата `t.me/joinchat/RWltlc4VqhWRzezx` или `@link`"
		err1 := errors.New(errMsg)
		return errMsg, err1
	}
	return "", nil
}
func botAwait(token string) {
	//	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	g := color.New(color.FgGreen, color.Bold).SprintFunc()
	b := color.New(color.FgBlue, color.Bold).SprintFunc()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("%s %s %s\n", red("error>"), g("botAwait>"), yellow(err))
		//	red.Printf("  %s\n")
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
				//if args == "" {
				//	msg.Text = "Не указана ссылка на чат. Требуется ссылка формата `t.me/joinchat/RWltlc4VqhWRzezx`"
				//}

				//------------------------------------------------
				/*
					var user []ClientConfig
					DB.Where("city = ? and login = ?", city, login).Find(&user)
					fmt.Println(user)
					if len(user) == 0 {
						if result := DB.Create(&conf); result.Error != nil {
							return nil, "", fmt.Errorf("conf %s with login %s is exists", login, chatLink)
						}
						return &conf, "create", nil
					}
					rconf := ClientConfig{
						Login:    login,
						City:     city,
						ChatLink: chatLink,
						BotToken: token,
					}
					DB.Model(&rconf).Where("city = ? and login = ?", city, login).Updates(rconf)
					//DB.Model(&ClientConfig{}).Where("login = ?", user[0].Login).Update("name", "hello")
					return &rconf, "update", nil
				*/
				//------------------------------------------------

				var links []Link
				DAL.DB.Where("chat_id = ?", message.Chat.ID).Find(&links)
				fmt.Printf("%s %s link: %s\n", red("botAwait>"), g("find link>"), yellow(link))
				link = Link{
					UserLink: args,
					ChatID:   message.Chat.ID,
				}
				if len(links) == 0 {
					msg.Text = "Чат еще не привязан. Привязываем."
					final, _ := bot.Send(msg)
					if result := DB.Create(&link); result.Error != nil {
						fmt.Printf("%s  %s  %d is exists\n", red("error>"), args, uint64(message.Chat.ID))
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

				//------------------------------------------------
				//DAL.DB.Where("chat_id = ?", message.Chat.ID).First(&link)
				///if len(link) != 0 {

				//}
				fmt.Printf("%s link: %s\n", g("botAwait>"), yellow(link))
				/*	if (Link{} == link) {
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
				*/
			}
		}

		//log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)
		fmt.Printf("%s %s %s %s\n", red("message>"), yellow(message.Chat.ID), b(message.From.UserName+">"), g(message.Text))
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
