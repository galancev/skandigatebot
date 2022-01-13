package bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	a "skandigatebot/models/account"
	u "skandigatebot/models/user"
)

func SendMessage(message string, m *tb.Message, b *tb.Bot) {
	_, err := b.Send(m.Sender, message, tb.ModeHTML)
	if err != nil {
		log.Fatal(err)
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
