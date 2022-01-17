package bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	a "skandigatebot/models/account"
	u "skandigatebot/models/user"
)

type LogChat struct {
	LogChatId string
}

func (lc LogChat) Recipient() string {
	return LogChatId
}

const (
	LogChatId = "-615741784"
)

func SendMessage(message string, m *tb.Message, b *tb.Bot) {
	_, err := b.Send(m.Sender, message, tb.ModeHTML)
	if err != nil {
		log.Print(err)
	}
}

func GetAccountAndUser(m *tb.Message) (a.Account, u.User, error) {
	account := a.GetAccount(m)
	var err error
	var user u.User

	if account.Phone > 0 {
		user, err = u.GetUser(account.Phone)
	}

	return account, user, err
}

func SendMessageLog(message string, b *tb.Bot) {
	var r LogChat

	r.LogChatId = LogChatId

	_, err := b.Send(r, message, tb.ModeHTML)
	if err != nil {
		log.Print(err)
	}
}
