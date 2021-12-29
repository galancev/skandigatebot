package gate

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"skandigatebot/bot"
	a "skandigatebot/models/account"
	u "skandigatebot/models/user"
	"skandigatebot/models/user/role"
	"skandigatebot/screens/admin"
	"time"
)

const (
	textSelectAction = "Выберите дальнейшее действие"
	OpenGateButton   = "Открыть врата!"
	textGateOpening  = "Врата открываются..."
	textNonAuth      = "Вам нельзя это сделать, вы не авторизованы."
)

type pauth interface {
	ShowAuthMenu(account *a.Account, user *u.User, m *tb.Message, b *tb.Bot)
}

type PGate struct {
	PAuth pauth
}

func New(pauth pauth) *PGate {
	return &PGate{
		PAuth: pauth,
	}
}

func (pg *PGate) OnOpen(m *tb.Message, b *tb.Bot) {
	account, user, _ := bot.GetAccountAndUser(m)

	if account.Phone > 0 && user.Phone > 0 {
		pg.HideGateMenuWithMessage(textGateOpening, &account, &user, m, b)

		OpenGate(m, b)

		pg.ShowGateMenu(&account, &user, m, b)
	} else {
		bot.SendMessage(textNonAuth, m, b)
		pg.PAuth.ShowAuthMenu(&account, &user, m, b)
	}
}

func (pg *PGate) ShowGateMenu(account *a.Account, user *u.User, m *tb.Message, b *tb.Bot) {
	pg.ShowGateMenuWithMessage(textSelectAction, account, user, m, b)
}

func (pg *PGate) ShowGateMenuWithMessage(message string, account *a.Account, user *u.User, m *tb.Message, b *tb.Bot) {
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnOpenGate := menu.Text(OpenGateButton)
	btnAdmin := menu.Text(admin.OnAdminButton)

	if user.RoleId == role.Admin {
		menu.Reply(
			menu.Row(btnOpenGate, btnAdmin),
		)
	} else {
		menu.Reply(
			menu.Row(btnOpenGate),
		)
	}

	_, err := b.Send(m.Sender, message, menu)
	if err != nil {
		log.Fatal(err)
	}
}

func (pg *PGate) HideGateMenu(account *a.Account, user *u.User, m *tb.Message, b *tb.Bot) {
	pg.HideGateMenuWithMessage(textSelectAction, account, user, m, b)
}

func (pg *PGate) HideGateMenuWithMessage(message string, account *a.Account, user *u.User, m *tb.Message, b *tb.Bot) {
	menu := &tb.ReplyMarkup{ReplyKeyboardRemove: true}

	_, err := b.Send(m.Sender, message, menu)
	if err != nil {
		log.Fatal(err)
	}
}

func OpenGate(m *tb.Message, b *tb.Bot) {
	time.Sleep(5 * time.Second)

	bot.SendMessage("Врата открыты!", m, b)
}