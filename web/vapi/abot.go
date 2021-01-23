package vapi

import (
	"fmt"

	"github.com/BurntSushi/toml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/web/DAL"
)

var (
	token = ""
	//C *SNBot
	C *SNBot
	//Skip bool
	Skip bool
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

	Updatetoken := make(chan string)
	var user []DAL.ClientConfig
	DB.Where("").Find(&user)
	var config tomlConfig
	if _, err := toml.DecodeFile("token.toml", &config); err != nil {
		fmt.Println(err)
	}
	//	token = config.Token
	s := &Config{
		Token:      config.Token,
		UpdateTime: 60,
	}

	C, _ = New(s)

	for {
		select {
		case tok := <-Updatetoken:

			if C.bot.Token == tok {
				fmc.Printfln("#rbtskip Token> #gbt%s", tok)
			} else {
				fmc.Printfln("#rbtchange Token> #gbt%s", tok)
				C.bot.StopReceivingUpdates()
				s := &Config{
					Token:      tok,
					UpdateTime: 60,
				}
				C, _ = New(s)
			}

		case update := <-C.upd:
			fmc.Printfln("#rbt message>#ybt %d>#btt %s> #gbt%s", update.Message.Chat.ID, update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			C.Send(update.Message.Chat.ID, update.Message.Text)
		}
	}

}

//Send Send(chatID int64, msg string) send Message to chat by id
func (s *SNBot) Send(chatID int64, msg string) {
	m := tgbotapi.NewMessage(chatID, msg)
	_, _ = s.bot.Send(m)
}
