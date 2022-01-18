package users

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"skandigatebot/bot"
	u "skandigatebot/models/user"
	"skandigatebot/models/user/role"
)

type PAdminUsers struct {
	PAuth pauth
	PGate pgate
}

func (pau *PAdminUsers) OnAdminUsers(m *tb.Message, b *tb.Bot) {
	account, user, err := bot.GetAccountAndUser(m)

	if account.Phone > 0 {
		if err != nil {
			if err == u.ErrNotFound {
				bot.SendMessage(textAuthAccessDenied, m, b)
			} else {
				bot.SendMessage(textDbError, m, b)
			}
			pau.PAuth.ShowAuthMenu(&account, &user, m, b)
		} else {
			if user.RoleId == role.Admin {
				pau.ShowUserList(m, b)
			} else {
				bot.SendMessage(textAuthAdminDenied, m, b)
				pau.PGate.ShowGateMenu(&account, &user, m, b)
			}
		}
	} else {
		bot.SendMessage(textNonAuth, m, b)
		pau.PAuth.ShowAuthMenu(&account, &user, m, b)
	}
}

func (pau *PAdminUsers) ShowUserList(m *tb.Message, b *tb.Bot) {
	var currentPage int
	currentPage = 1

	selector := getAdminUserSelector(currentPage, m, b)

	_, err := b.Send(m.Sender, getAdminUserMessage(currentPage), selector, tb.ModeHTML)
	if err != nil {
		bot.SendMessageLog(err.Error(), b)
	}
}
