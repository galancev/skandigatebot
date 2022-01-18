package admin

import (
	tb "gopkg.in/tucnak/telebot.v2"
	a "skandigatebot/models/account"
	u "skandigatebot/models/user"
)

type pauth interface {
	ShowAuthMenu(account *a.Account, user *u.User, m *tb.Message, b *tb.Bot)
}
