package admin

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"skandigatebot/bot"
	u "skandigatebot/models/user"
	"skandigatebot/models/user/role"
)

const (
	textSelectAction     = "🤔 Выберите дальнейшее действие"
	OnAdminButton        = "😇 Админка"
	textDbError          = "😵 Проблема с базой данных на сервере"
	textAuthAccessDenied = "❗️ Вы успешно авторизовались, однако вашего телефона нет в списке разрешенных. По вопросам доступа обратитесь в управляющую организацию АО «ВК Комфорт». Доступ возможен только для разгрузки и проезда на смежные территории."
	textAuthAdminDenied  = "📛 Хорошая попытка, но нет. В админку вам нельзя!"
	textNonAuth          = "⛔️ Вам нельзя это сделать, вы не авторизованы."
	OnAdminExitButton    = "↩️ Выйти из админки"
	OnAdminShowUsers     = "👥 Показать пользователей"
)

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
		bot.SendMessageLog(err.Error(), b)
	}
}
