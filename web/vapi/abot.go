package vapi

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/web/DAL"
)

type MessageTG struct {
	Message string
	ChatID  string
}

var (
	token = ""
	//C *SNBot
	C *SNBot
	//Skip bool
	Skip             bool
	MessageTGChannel chan MessageTG
	Updatetoken      chan string
)

//New (cfg *Config) (*SNBot, error)
func New(cfg *Config) (*SNBot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.UpdateTime
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return &SNBot{
		cfg: cfg,
		bot: bot,
		upd: updates,
	}, nil
}

func runBot() {
	fmc.Println("#rbtrun bot> run")

	var user []DAL.ClientConfig
	DB.Where("").Find(&user)
	//var config tomlConfig
	//if _, err := toml.DecodeFile("token.toml", &config); err != nil {
	//	fmt.Println(err)
	//}
	//	token = config.Token
	s := &Config{
		Token:      user[0].BotToken,
		UpdateTime: 60,
	}

	C, _ = New(s)
	for {

		select {

		case tok := <-Updatetoken:
			var user []DAL.ClientConfig
			DB.Where("").Find(&user)
			fmc.Printfln("user:%s", user[0].BotToken)
			if C.bot.Token == user[0].BotToken {
				fmc.Printfln("#rbtskip Token> #gbt%s", tok)
			} else {
				fmc.Printfln("#rbtchange Token> #gbt%s", user[0].BotToken)
				C.bot.StopReceivingUpdates()
				s := &Config{
					Token:      user[0].BotToken,
					UpdateTime: 60,
				}
				C, _ = New(s)
			}

		case update := <-C.upd:
			switch update.Message.Command() {
			case "id":
				msg := fmt.Sprintln("ID Чата: " + strconv.FormatInt(update.Message.Chat.ID, 10))
				C.Send(update.Message.Chat.ID, msg)
			case "check":
				var link Link
				msg := ""
				//	msg.ParseMode = tgbotapi.ModeMarkdown
				DAL.DB.Where("chat_id = ?", update.Message.Chat.ID).First(&link)
				if (link != Link{}) {
					msg = fmt.Sprintf("Чат привязан. Сводка: \n *ChatID*: %d \n UserLink: %s", link.ChatID, link.UserLink)
				} else {
					msg = fmt.Sprintf("Чат *не привязан*")
				}
				C.Send(update.Message.Chat.ID, msg)
			case "link":
				admins, err := C.bot.GetChatAdministrators(update.Message.Chat.ChatConfig())
				if err != nil {
					msg := "Вероятно что-то пошло не так. Проверьте права админа у бота."
					C.Send(update.Message.Chat.ID, msg)
					//	return
				} else {
					msg := fmt.Sprint(admins)
					C.Send(update.Message.Chat.ID, msg)

				}
				//if !userIsAdmin(update.Message.From, admins) {
				//	msg.Text = "У тебя нет прав администратора."

				//	return
				//	}
				//	errMsg, errLink := CheckLink(args)
				//	if errLink != nil {
				//		msg.Text = errMsg
				//		bot.Send(msg)
				//		return
				//	}
			}
			fmc.Printfln("#rbt message>#ybt %d>#btt %s> #gbt%s", update.Message.Chat.ID, update.Message.From.UserName, update.Message.Text)
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID
			//C.Send(update.Message.Chat.ID, update.Message.Text)
		case msg := <-MessageTGChannel:
			var links []Link
			DB.Where("user_link = ?", msg.ChatID).Find(&links)
			fmt.Println(links[0])
		}
	}

}

//Send Send(chatID int64, msg string) send Message to chat by id
func (s *SNBot) Send(chatID int64, msg string) {
	m := tgbotapi.NewMessage(chatID, msg)
	_, _ = s.bot.Send(m)
}
