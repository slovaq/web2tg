package vapi

import (
	"fmt"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/slovaq/web2tg/web/DAL"
)

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	g      = color.New(color.FgGreen, color.Bold).SprintFunc()
	b      = color.New(color.FgBlue, color.Bold).SprintFunc()
	token  = ""
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

func (Update Update) edit(w http.ResponseWriter, r *http.Request) {
	token = r.FormValue("token")
	//fmt.Printf("%s %s %s\n", yellow("edit>"), b("token>"), g(token))
	w.Write([]byte(token))
	go func() {
		Update.UpdateToken <- token

	}()

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
				fmt.Printf("%s %s\n", red("skip Token>"), g(tok))
			} else {
				fmt.Printf("%s %s\n", b("change Token>"), g(tok))
				C.bot.StopReceivingUpdates()
				s := &Config{
					Token:      tok,
					UpdateTime: 60,
				}
				C, _ = New(s)
			}

		case update := <-C.upd:
			fmt.Printf("%s %s %s %s\n", red("message>"), yellow(update.Message.Chat.ID, ">"), b(update.Message.From.UserName+">"), g(update.Message.Text))
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
