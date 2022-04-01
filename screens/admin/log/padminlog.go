package log

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"skandigatebot/bot"
	u "skandigatebot/models/user"
)

type PAdminLog struct {
	PAuth pauth
	PGate pgate
}

func New(pauth pauth, pgate pgate) *PAdminLog {
	return &PAdminLog{
		PAuth: pauth,
		PGate: pgate,
	}
}

func (pal *PAdminLog) OnAdminLog(m *tb.Message, b *tb.Bot) {
	account, user, err := bot.GetAccountAndUser(m)

	if account.Phone > 0 {
		if err != nil {
			if err == u.ErrNotFound {
				bot.SendMessage(textAuthAccessDenied, m, b)
			} else {
				bot.SendMessage(textDbError, m, b)
			}
			pal.PAuth.ShowAuthMenu(&account, &user, m, b)
		} else {
			if user.IsBlocked() {
				bot.SendMessage(textAuthAccessDenied, m, b)
				pal.PAuth.ShowAuthMenu(&account, &user, m, b)
			} else {
				if user.IsAdmin() {
					pal.ShowLog(m, b)
				} else {
					bot.SendMessage(textAuthAdminDenied, m, b)
					pal.PGate.ShowGateMenu(&account, &user, m, b)
				}
			}
		}
	} else {
		bot.SendMessage(textNonAuth, m, b)
		pal.PAuth.ShowAuthMenu(&account, &user, m, b)
	}
}

func (pau *PAdminLog) ShowLog(m *tb.Message, b *tb.Bot) {
	var currentPage int
	currentPage = 1

	selector := getLogUserSelector(currentPage, m, b)

	_, err := b.Send(m.Sender, getAdminLogMessage(currentPage), selector, tb.ModeHTML)
	if err != nil {
		bot.SendMessageLog(err.Error(), b)
	}
}
