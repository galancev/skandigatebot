package user

type PAdminUser struct {
	/*PAuth pauth
	PGate pgate*/
}

func New( /*pauth pauth, pgate pgate*/ ) *PAdminUser {
	return &PAdminUser{
		/*PAuth: pauth,
		PGate: pgate,*/
	}
}

/*func (pau *PAdminUser) OnAdminUsers(m *tb.Message, b *tb.Bot) {
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

func (pau *PAdminUser) ShowUser(userId int, page int, m *tb.Message, b *tb.Bot) {
	selector := getAdminUserSelector(currentPage, m, b)

	_, err := b.Send(m.Sender, getAdminUserMessage(currentPage), selector, tb.ModeHTML)
	if err != nil {
		bot.SendMessageLog(err.Error(), b)
	}
}*/
