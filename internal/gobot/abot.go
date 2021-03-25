package gobot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/internal/DAL"
)

var (
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
func TestRC(upd *GobotConnect) {
	fmc.Printfln("#rbt Run> #gbtInitBotRC")
	//upd.InitBot()
}
func (upd *GobotConnect) InitBot() {
	fmc.Println("#rbtInitBot> #gbtstart")
	var user []DAL.ClientConfig
	DAL.DB.Where("").Find(&user)
	initV := 0
	if len(user) != 0 {
		//fmt.Println("init>user>", user[0].BotToken)
		//fmc.Println("#rbtinitBot> Run bot>")
		go upd.RunBot()
		initV = 1
	} else {
		fmc.Println("#rbtinitBot>False Run>")
	}
	for range upd.CheckInit {
		if initV == 0 {
			//fmc.Println("#rbtinitBot> Run bot>")
			go upd.RunBot()
			initV = 1
		} else {
			fmc.Println("#rbtinitBot>Update Token>")
			upd.Updatetoken <- true
		}
	}
}
func (upd *GobotConnect) RunBot() {
	fmc.Println("#ybtRunBot> #gbtstart")
	var user []DAL.ClientConfig
	DAL.DB.Where("").Find(&user)
	//for _, u := range user {
	u := user[0]
	//go func() {
	s := &Config{
		Token:      u.BotToken,
		UpdateTime: 60,
	}

	C, err := New(s)
	if err != nil {
		fmc.Printfln("err: %s", err)
		//select {}
	}
	C.bot.Debug = false
	log.Printf("Authorized on account %s", C.bot.Self.UserName)
	fmc.Printfln("#ybtrange> #gbtstart")
	//for x := range upd.MessageTG {

	//	fmc.Printfln("range> %s", x.ChatID)
	//}
	for {
		//fmc.Println("#ybt(loop) #gbtselect")
		select {
		case <-upd.Updatetoken:
			var user []DAL.ClientConfig
			DAL.DB.Where("").Find(&user)
			for _, u := range user {
				fmc.Printfln("user:%s", u.BotToken)
				if C.bot.Token == u.BotToken {
					fmc.Printfln("#rbtskip Token> #gbt%s", u.BotToken)
				} else {
					fmc.Printfln("#rbtchange Token> #gbt%s", u.BotToken)
					C.bot.StopReceivingUpdates()
					s := &Config{
						Token:      u.BotToken,
						UpdateTime: 60,
					}
					C, _ = New(s)
				}
			}
		case update := <-C.upd:

			if update.EditedMessage != nil {
				fmc.Printfln("#ybtedited message> #gbtmsg: %s", update.EditedMessage.Text)

			} else {
				fmc.Printfln("#ybtmsg: #bbt%s %s> #gbt %s", update.Message.From.FirstName, update.Message.From.LastName, update.Message.Text)
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "id":
						returnChatid(C.bot, update.Message)
					case "check":
						if update.Message.CommandArguments() == "" {
							C.Send(update.Message.Chat.ID, "/check https://t.me/joinchat/asdhajksdasd")
						} else {
							id, msg := checkChat(update.Message, update.Message.CommandArguments()) // По другому никак.
							C.Send(id, msg)
							fmc.Printfln("#gbtid: %d, #bbtupdate.id: %d, msg:%s", id, update.Message.Chat.ID, msg)

						}

					case "link":
						linkChat(C, update.Message)
					default:
						fmc.Printfln("#rbtCommandHandler Error> command not found: #gbt%s", update.Message.Command())
					}
				}
			}

		case msg := <-upd.MessageTG:
			fmc.Println("#rbtupd.MessageTG start> ")
			var links Link
			DAL.DB.Where("user_link = ?", msg.ChatID).First(&links)
			if (links == Link{}) {
				fmc.Println("#rbtupd.MessageTG Error> #ybtlen(links) = 0")
				return
			}
			ms, err := C.Send(links.ChatID, msg.Message)

			if msg.Pic != "" {
				ms, err = C.SendPhoto(links.ChatID, msg.Pic)
			}

			fmc.Printfln("<-MessageTGChannel : %d, %s", links.ChatID, msg.Message)
			fmc.Printfln("#ybtpic> #gbt%s", msg.Pic)

			if err != nil {
				fmc.Printfln("message error>%s \n\t ms: %v", err.Error(), ms)
			}
			fmc.Println("#rbtupd.MessageTG stop> ")
		}
	}
	//	}()
	//}
}

//Send Send(chatID int64, msg string) send Message to chat by id
func (s *SNBot) Send(chatID int64, msg string) (tgbotapi.Message, error) {
	m := tgbotapi.NewMessage(chatID, msg)
	m.ParseMode = tgbotapi.ModeMarkdown
	return s.bot.Send(m)
}

//SendPhoto Send(chatID int64, msg string) send Message to chat by id
func (s *SNBot) SendPhoto(chatID int64, pic string) (tgbotapi.Message, error) {
	m := tgbotapi.NewPhotoUpload(chatID, pic)
	return s.bot.Send(m)
}
