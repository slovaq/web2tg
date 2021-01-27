package vapi

import (
	"fmt"
	"log"

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

func (upd *UpdateStorage) runBot() {
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

	C, err := New(s)
	C.bot.Debug = true
	log.Printf("Authorized on account %s", C.bot.Self.UserName)

	if err != nil {

		fmc.Printfln("err: %s", err)
		for {
		}
	}
	for {
		select {

		case <-upd.Updatetoken:
			var user []DAL.ClientConfig
			DB.Where("").Find(&user)
			fmc.Printfln("user:%s", user[0].BotToken)
			if C.bot.Token == user[0].BotToken {
				fmc.Printfln("#rbtskip Token> #gbt%s", user[0].BotToken)
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
			fmt.Println(update.Message)
			//time.Sleep(1 * time.Second)
			if update.Message.IsCommand() == true {

				switch update.Message.Command() {
				case "id":
					returnChatid(C.bot, update.Message)
				case "check":
					id, msg := checkChat(C, update.Message, user[0].ChatLink)
					C.Send(id, msg)
					fmc.Printfln("#gbtid: %d, #bbtupdate.id: %d, msg:%s", id, update.Message.Chat.ID, msg)
				case "link":
					linkChat(C, update.Message)
				}
			}
		//
		case msg := <-MessageTGChannel:
			var links []Link
			DB.Where("user_link = ?", msg.ChatID).Find(&links)
			fmt.Println(links[0])
		}
	}

}

//Send Send(chatID int64, msg string) send Message to chat by id
func (s *SNBot) Send(chatID int64, msg string) (tgbotapi.Message, error) {
	m := tgbotapi.NewMessage(chatID, msg)
	m.ParseMode = tgbotapi.ModeHTML
	return s.bot.Send(m)
}
