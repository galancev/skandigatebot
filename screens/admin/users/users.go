package users

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"skandigatebot/bot"
	a "skandigatebot/models/account"
	u "skandigatebot/models/user"
	"skandigatebot/models/user/role"
	"strconv"
)

const (
	textDbError          = "Конина какая-то на сервере"
	textAuthAccessDenied = "Вы успешно авторизовались, однако вашего телефона нет в списке разрешённых. Напишите скандифокс для добавления."
	textAuthAdminDenied  = "Хорошая попытка, но нет. В админку вам нельзя!"
	textNonAuth          = "Вам нельзя это сделать, вы не авторизованы."

	userPerPage = 10
)

type pauth interface {
	ShowAuthMenu(account *a.Account, user *u.User, m *tb.Message, b *tb.Bot)
}

type pgate interface {
	ShowGateMenu(account *a.Account, user *u.User, m *tb.Message, b *tb.Bot)
}

type PAdminUsers struct {
	PAuth pauth
	PGate pgate
}

func New(pauth pauth, pgate pgate) *PAdminUsers {
	return &PAdminUsers{
		PAuth: pauth,
		PGate: pgate,
	}
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

	selector := &tb.ReplyMarkup{}

	btnPrev := selector.Data("⬅", "prev")
	btnNext := selector.Data("➡", "next")

	selector.Inline(
		selector.Row(btnPrev, btnNext),
	)

	send, err := b.Send(m.Sender, getAdminUserMessage(currentPage), selector)
	if err != nil {
		return
	}

	b.Handle(&btnPrev, func(c *tb.Callback) {
		usersCount, _ := u.GetUsersCount()
		pagesCount := usersCount/userPerPage + 1

		currentPage--

		if currentPage < 1 {
			currentPage = int(pagesCount)
		}

		b.Edit(send, getAdminUserMessage(currentPage), selector)

		// ...
		// Always respond!
		b.Respond(c, &tb.CallbackResponse{})
	})

	b.Handle(&btnNext, func(c *tb.Callback) {
		usersCount, _ := u.GetUsersCount()
		pagesCount := usersCount/userPerPage + 1

		currentPage++

		if currentPage > int(pagesCount) {
			currentPage = 1
		}

		b.Edit(send, getAdminUserMessage(currentPage), selector)

		// ...
		// Always respond!
		b.Respond(c, &tb.CallbackResponse{})
	})
}

func getAdminUserMessage(page int) string {
	usersCount, _ := u.GetUsersCount()
	pagesCount := usersCount/userPerPage + 1

	users, err := u.GetUsers((page-1)*userPerPage, userPerPage)
	if err != nil {
		return ""
	}

	var message string

	for _, user := range users {
		message += "+" + strconv.Itoa(int(user.Phone))
		message += " " + user.UserFirstName + " " + user.UserLastName
		message += "\n"

		if user.AccountFirstName != "" {
			message += user.AccountFirstName + " " + user.AccountLastName

			if user.AccountUserName != "" {
				message += " @" + user.AccountUserName
			}

			message += "\n"
		}

		if user.RoleId == role.Admin {
			message += " [admin]\n"
		}

		message += "\n"
	}

	message += "[" + strconv.Itoa(page) + "] из [" + strconv.Itoa(int(pagesCount)) + "]"

	return message
}
