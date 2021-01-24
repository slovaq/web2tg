package vapi

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/slovaq/web2tg/web/DAL"
)

func userIsAdmin(member *tgbotapi.User, members []tgbotapi.ChatMember) bool {
	for _, admin := range members {
		if member == *admin.User {
			return true
		}
	}
	return false
}
func returnChatid(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "ID Чата"+strconv.FormatInt(message.Chat.ID, 10))
	bot.Send(msg)
}
func checkChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
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

func linkChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var link Link
	msg := tgbotapi.NewMessage(message.Chat.ID, "Линка привязана. ")
	msg.ParseMode = tgbotapi.ModeMarkdown
	args := message.CommandArguments()
	admins, err := bot.GetChatAdministrators(message.Chat.ChatConfig())
	if err != nil {
		msg.Text = "Вероятно что-то пошло не так. Проверьте права админа у бота."
		bot.Send(msg)
		return
	}
	if !userIsAdmin(*message.From, admins) {
		msg.Text = "У тебя нет прав администратора."
		bot.Send(msg)
		return
	}
	errMsg, errLink := CheckLink(args)
	if errLink != nil {
		msg.Text = errMsg
		bot.Send(msg)
		return
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
		return
	}

	DB.Model(&link).Where("chat_id = ?", message.Chat.ID).Updates(link)
	msg.Text = "Ссылка обновлена! Проверить можно по команде /check"
	bot.Send(msg)
	return

}
