package main

import (
	"fmt"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Config struct {
	Token      string
	UpdateTime int
}

type SNBot struct {
	cfg *Config
	bot *tgbotapi.BotAPI
	upd tgbotapi.UpdatesChannel
}

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

type tomlConfig struct {
	Token string
}

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	g      = color.New(color.FgGreen, color.Bold).SprintFunc()
	b      = color.New(color.FgBlue, color.Bold).SprintFunc()
)
var token string

type Update struct {
	UpdateToken chan string
}

func (Update Update) edit(w http.ResponseWriter, r *http.Request) {
	token = r.FormValue("token")
	fmt.Printf("%s %s %s\n", yellow("edit>"), b("token>"), g(token))
	w.Write([]byte(token))
	go func() {
		Update.UpdateToken <- token

	}()

}

var C *SNBot
var Skip bool

func main() {
	Updatetoken := make(chan string)
	//go
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	Upd := Update{

		UpdateToken: Updatetoken,
	}
	r.Get("/edit", Upd.edit)
	//curl "http://localhost:3000/edit?token=713753713:AAHvUVHW9MLQ1OVdXzumghRXOj_lShalCfQ"
	//curl "http://localhost:3000/edit?token=697389856:AAFxgMjR6yMMjHek1KfNaXrikRYNkZmuJww"
	fmt.Println("listen:3000")
	go http.ListenAndServe(":3000", r)
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

	//	Updatetoken := make(chan string)
	//memoryLatest := 0
	//latestMessage := 0
	for {
		select {
		case tok := <-Updatetoken:
			//	fmt.Printf("%s %s\n", yellow("UpdateToken:"), g(tok))
			//fmt.Println(latestMessage)
			//	hook, err := C.bot.RemoveWebhook()
			//if err != nil {
			//		fmt.Printf("could not remove webhook: %v", err)
			//	}
			//fmt.Println("token: ", C.bot.Token)
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

		//	hook2, err := C.bot.RemoveWebhook()
		//	if err != nil {
		//		fmt.Printf("could not remove webhook: %v", err)
		//	}
		//fmt.Println("hook: ", hook2)
		//Skip = true
		case update := <-C.upd:
			//if Skip == true {
			//	if memoryLatest < update.Message.MessageID {
			//	fmt.Printf("memoryLatest: %d, update.Message.MessageID:%d\n ", memoryLatest, update.Message.MessageID)
			//
			//	} else {

			//Skip = false
			//	}
			//latestMessage = update.Message.MessageID
			//} else {

			//	latestMessage = update.Message.MessageID
			//fmt.Println("quit")
			fmt.Printf("%s %s %s %s\n", red("message>"), yellow(update.Message.Chat.ID, ">"), b(update.Message.From.UserName+">"), g(update.Message.Text))
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			C.Send(update.Message.Chat.ID, update.Message.Text)
			//	}
			//
		}
	}

}
func (s *SNBot) Send(chatID int64, msg string) {
	m := tgbotapi.NewMessage(chatID, msg)
	_, _ = s.bot.Send(m)
}
