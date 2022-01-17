package auth

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"skandigatebot/base"
	"skandigatebot/bot"
	a "skandigatebot/models/account"
	u "skandigatebot/models/user"
	"skandigatebot/screens/first"
	"skandigatebot/screens/gate"
	"strconv"
)

const (
	textSelectAction     = "🤔 Выберите дальнейшее действие"
	textEnemyPhoneNumber = "🤦‍♂️ Необходимо делиться своим телефоном, а не чужим!"
	textSharePhoneNumber = "📱 Поделиться номером телефона"
)

type PAuth struct{}

func New() *PAuth {
	return &PAuth{}
}

func (pa *PAuth) OnAuth(m *tb.Message, b *tb.Bot) {
	var user u.User

	account := a.GetAccount(m)

	if account.AccountId != uint(m.Contact.UserID) {
		bot.SendMessage(textEnemyPhoneNumber, m, b)

		return
	}

	phone, err := strconv.Atoi(m.Contact.PhoneNumber)
	if err != nil {
		bot.SendMessageLog(err.Error(), b)
	}

	if account.Phone != uint(phone) {
		account.Phone = uint(phone)

		base.GetDB().Save(&account)
	}

	account, user, _ = bot.GetAccountAndUser(m)

	pgate := gate.New(pa)
	pfirst := first.New(pa, pgate)
	pfirst.ShowFirstMenu(&account, &user, m, b)
}

func (pa *PAuth) ShowAuthMenu(account *a.Account, user *u.User, m *tb.Message, b *tb.Bot) {
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnSharePhoneNumber := menu.Contact(textSharePhoneNumber)

	menu.Reply(
		menu.Row(btnSharePhoneNumber),
	)

	_, err := b.Send(m.Sender, textSelectAction, menu)
	if err != nil {
		bot.SendMessageLog(err.Error(), b)
	}
}
