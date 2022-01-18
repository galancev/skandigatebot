package first

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"skandigatebot/bot"
	a "skandigatebot/models/account"
	u "skandigatebot/models/user"
)

const (
	textHello               = "🖖 Привет! Это скандибот для управления шлагбаумом паркинга 1 корпуса."
	textNeedAuth            = "‼️ Для продолжения работы необходимо авторизоваться."
	textAlreadyAuth         = "✅ Вы авторизованы, можете пользоваться шлагбаумом."
	textAuthButAccessDenied = "❗️ Вы успешно авторизовались, однако вашего телефона нет в списке разрешённых. Напишите @ScandiFox для добавления."
	textDbError             = "😵 Конина какая-то на сервере"
)

type PFirst struct {
	PAuth pauth
	PGate pgate
}

func New(pauth pauth, pgate pgate) *PFirst {
	return &PFirst{
		PAuth: pauth,
		PGate: pgate,
	}
}

func (pf *PFirst) OnStart(m *tb.Message, b *tb.Bot) {
	bot.SendMessage(textHello, m, b)

	account, user, _ := bot.GetAccountAndUser(m)

	pf.ShowFirstMenu(&account, &user, m, b)
}

func (pf *PFirst) ShowFirstMenu(account *a.Account, user *u.User, m *tb.Message, b *tb.Bot) {
	if account.Phone > 0 {
		var err error

		*user, err = u.GetUser(account.Phone)

		if err != nil {
			if err == u.ErrNotFound {
				bot.SendMessage(textAuthButAccessDenied, m, b)
			} else {
				bot.SendMessage(textDbError, m, b)
			}
			pf.PAuth.ShowAuthMenu(account, user, m, b)

		} else {
			bot.SendMessage(textAlreadyAuth, m, b)
			pf.PGate.ShowGateMenu(account, user, m, b)
		}
	} else {
		bot.SendMessage(textNeedAuth, m, b)
		pf.PAuth.ShowAuthMenu(account, user, m, b)
	}
}
