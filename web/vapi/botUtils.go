package vapi

import (
	"fmt"
	"github.com/mallvielfrass/coloredPrint/fmc"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/slovaq/web2tg/web/DAL"
)

func userIsAdmin(member *tgbotapi.User, members []tgbotapi.ChatMember) bool {
	for _, admin := range members {
		if member == admin.User {
			return true
		}
	}
	return false
}
func returnChatid(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "ID Чата: "+strconv.FormatInt(message.Chat.ID, 10))
	bot.Send(msg)
}
func checkChat(bot *SNBot, message *tgbotapi.Message) {
	var link Link
	var msg string
	DAL.DB.Where("chat_id = ?", message.Chat.ID).First(&link)
	if (link != Link{}) {
		msg = fmt.Sprintf("Чат привязан. Сводка: \n *ChatID*: %d \n UserLink: %s", link.ChatID, link.UserLink)
	} else {
		msg = fmt.Sprintf("Чат *не привязан*")
	}
	bot.Send(message.Chat.ID, msg)

}

func linkChat(bot *SNBot, message *tgbotapi.Message) {
	var link Link
	var msg string
	id := message.Chat.ID
	args := message.CommandArguments()
	// а чо оно не работает
	//if userIsAdmin(message.From, admins) {
	//	msg = "У тебя нет прав администратора."
	//	bot.Send(id, msg)
	//	return
	//}
	errMsg, errLink := CheckLink(args)
	if errLink != nil {
		msg = errMsg
		bot.Send(id, msg)
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
		msg = "Чат еще не привязан. Привязываем."
		final, _ := bot.Send(id, msg)
		if result := DB.Create(&link); result.Error != nil {
			fmt.Printf("%s  %s  %d is exists\n", redPrint("error>"), args, uint64(message.Chat.ID))
			msg = "Произошла ошибка при привязке. Попробуйте еще раз или обратитесь к системному администратору!"
		} else {
			bot.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
				ChannelUsername: message.Chat.UserName,
				ChatID:          message.Chat.ID,
				MessageID:       final.MessageID,
			})
			msg = "Чат теперь привязан! Проверить можно по команде /check"
		}
		bot.Send(id, msg)
		return
	}

	DB.Model(&link).Where("chat_id = ?", message.Chat.ID).Updates(link)
	msg = "Ссылка обновлена! Проверить можно по команде /check"
	fmc.Printfln("#rbt message>#ybt %d>#btt %s> #gbt%s", message.Chat.ID, message.From.UserName, message.Text)

	bot.Send(id, msg)
	return

}
