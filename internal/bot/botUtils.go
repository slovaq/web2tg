package bot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mallvielfrass/coloredPrint/fmc"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/slovaq/web2tg/internal/DAL"
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
func checkChat(message *tgbotapi.Message, UserURLlink string) (int64, string) {
	var link Link
	var user ClientConfig
	var msg string

	DAL.DB.Where("chat_link = ? ", UserURLlink).First(&user)

	IsLinked := strconv.FormatBool(user == ClientConfig{})

	DAL.DB.Where("chat_id = ?", message.Chat.ID).Find(&link)

	fmt.Println("links: ", link)
	if link.UserLink != "" {

		msg = fmt.Sprintf("Чат привязан.\nЧат привязан к юзеру?: %s \nChatID: %d\nUserLink: %s", IsLinked, message.Chat.ID, link.UserLink)
	} else {
		msg = "Чат не привязан"
	}
	fmc.Printfln("#rbtcheck chat>#ybt msg> #bbt%s", msg)
	//bot.Send(message.Chat.ID, msg)
	return message.Chat.ID, msg

}

//CheckLink (msg string) (string, error)
func CheckLink(msg string) (string, error) {

	fmc.Println("#ybtCheckLink> #gbtmsg")
	if !(strings.Contains(msg, "@")) && !(strings.Contains(msg, "t.me")) || msg == "" {
		errmsg := "Не указана ссылка на чат. Требуется ссылка формата `t.me/joinchat/RWltlc4VqhWRzezx` или `@link`"
		err1 := errors.New("the link to the chat was not specified")
		return errmsg, err1

	}
	return "", nil
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
	fmc.Printf("#rbtbotAwait> #gbtfind link> link: %s\n", link.UserLink)
	link = Link{
		UserLink: args,
		ChatID:   message.Chat.ID,
	}

	if len(links) == 0 {
		msg = "Чат еще не привязан. Привязываем."
		final, _ := bot.Send(id, msg)
		if result := DAL.DB.Create(&link); result.Error != nil {
			fmc.Printfln("#rbterror> #ybt  %d is exists\n", args, uint64(message.Chat.ID))
			msg = "Произошла ошибка при привязке. Попробуйте еще раз или обратитесь к системному администратору!"
		} else {
			bot.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
				//ChannelUsername: message.Chat.UserName,
				ChatID:    message.Chat.ID,
				MessageID: final.MessageID,
			})
			msg = "Чат теперь привязан! Проверить можно по команде /check"
		}
		bot.Send(id, msg)
		return
	}

	DAL.DB.Model(&link).Where("chat_id = ?", message.Chat.ID).Updates(link)
	msg = "Ссылка обновлена! Проверить можно по команде /check"
	fmc.Printfln("#rbt message>#ybt %d>#btt %s> #gbt%s", message.Chat.ID, message.From.UserName, message.Text)

	bot.Send(id, msg)

}
