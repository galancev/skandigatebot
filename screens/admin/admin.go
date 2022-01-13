package admin

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"skandigatebot/bot"
	a "skandigatebot/models/account"
	u "skandigatebot/models/user"
	"skandigatebot/models/user/role"
)

const (
	textSelectAction     = "Выберите дальнейшее действие"
	OnAdminButton        = "Админка"
	textDbError          = "Конина какая-то на сервере"
	textAuthAccessDenied = "Вы успешно авторизовались, однако вашего телефона нет в списке разрешённых. Напишите скандифокс для добавления."
	textAuthAdminDenied  = "Хорошая попытка, но нет. В админку вам нельзя!"
	textNonAuth          = "Вам нельзя это сделать, вы не авторизованы."
	OnAdminExitButton    = "Выйти из админки"
	OnAdminShowUsers     = "Показать пользователей"
)

type pauth interface {
	ShowAuthMenu(account *a.Account, user *u.User, m *tb.Message, b *tb.Bot)
}

type pgate interface {
	ShowGateMenu(account *a.Account, user *u.User, m *tb.Message, b *tb.Bot)
}

type PAdmin struct {
	PAuth pauth
	PGate pgate
}

func New(pauth pauth, pgate pgate) *PAdmin {
	return &PAdmin{
		PAuth: pauth,
		PGate: pgate,
	}
}

func (pa *PAdmin) OnAdmin(m *tb.Message, b *tb.Bot) {
	account, user, err := bot.GetAccountAndUser(m)

	if account.Phone > 0 {
		if err != nil {
			if err == u.ErrNotFound {
				bot.SendMessage(textAuthAccessDenied, m, b)
			} else {
				bot.SendMessage(textDbError, m, b)
			}
			pa.PAuth.ShowAuthMenu(&account, &user, m, b)
		} else {
			if user.RoleId == role.Admin {
				pa.ShowAdminMenu(m, b)
			} else {
				bot.SendMessage(textAuthAdminDenied, m, b)
				pa.PGate.ShowGateMenu(&account, &user, m, b)
			}
		}
	} else {
		bot.SendMessage(textNonAuth, m, b)
		pa.PAuth.ShowAuthMenu(&account, &user, m, b)
	}
}

func (pa *PAdmin) ShowAdminMenu(m *tb.Message, b *tb.Bot) {
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnAdminBack := menu.Text(OnAdminExitButton)
	btnAdminShowUsers := menu.Text(OnAdminShowUsers)

	menu.Reply(
		menu.Row(btnAdminBack, btnAdminShowUsers),
	)

	_, err := b.Send(m.Sender, textSelectAction, menu)
	if err != nil {
		log.Fatal(err)
	}
}
