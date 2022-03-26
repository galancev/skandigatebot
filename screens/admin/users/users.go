package users

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"skandigatebot/bot"
	"skandigatebot/models"
	u "skandigatebot/models/user"
	"skandigatebot/models/user/role"
	"strconv"
	//au "skandigatebot/screens/admin/user"
)

const (
	textDbError          = "😵 Проблема с базой данных на сервере"
	textAuthAccessDenied = "❗️ Вы успешно авторизовались, однако вашего телефона нет в списке разрешенных. По вопросам доступа обратитесь в управляющую организацию АО «ВК Комфорт». Доступ возможен только для разгрузки и проезда на смежные территории."
	textAuthAdminDenied  = "📛 Хорошая попытка, но нет. В админку вам нельзя!"
	textNonAuth          = "⛔️ Вам нельзя это сделать, вы не авторизованы."

	userPerPage = 10
)

func getAdminUsers(page int) ([]models.UserAccount, error) {
	users, err := u.GetUsersWithAccount((page-1)*userPerPage, userPerPage)

	return users, err
}

func getAdminUserMessage(page int) string {
	usersCount, _ := u.GetUsersCount()
	pagesCount := usersCount/userPerPage + 1

	var message string

	message += "Страница [" + strconv.Itoa(page) + "] из [" + strconv.Itoa(int(pagesCount)) + "]"

	return message
}

func getAdminUserSelector(page int, m *tb.Message, b *tb.Bot) *tb.ReplyMarkup {
	selector := &tb.ReplyMarkup{}

	btnPrev := selector.Data("⬅", "prev", strconv.Itoa(page))
	btnNext := selector.Data("➡", "next", strconv.Itoa(page))

	users, err := getAdminUsers(page)
	if err != nil {
		return nil
	}

	var userButtons []tb.Btn
	for _, user := range users {
		message := ""
		if user.RoleId == role.Admin {
			message += "😇"
		} else {
			message += "👤"
		}

		message += "+" + strconv.Itoa(int(user.Phone))
		message += " (" + user.UserFirstName + ")"
		if user.AccountUserName != "" {
			message += " :: " + user.AccountFirstName + " " + user.AccountLastName
		}

		if user.AccountLastName != "" {
			message += " [" + user.AccountUserName + "]"
		}

		message += "\n"
		/*
			if user.AccountFirstName != "" {
				message += user.AccountFirstName + " " + user.AccountLastName

				if user.AccountUserName != "" {
					message += " @" + user.AccountUserName
				}
			}*/

		userButton := selector.Data(message, "account-"+strconv.Itoa(int(user.UserId)), strconv.Itoa(int(user.UserId)))

		userButtons = append(userButtons, userButton)

		b.Handle(&userButton, func(c *tb.Callback) {
			//pau := au.New()

			//pau.OnAdminUsers(m, b)

			err := b.Respond(c, &tb.CallbackResponse{})
			if err != nil {
				bot.SendMessageLog(err.Error(), b)
			}
		})
	}

	var rows []tb.Row

	for _, userButton := range userButtons {
		rows = append(rows, selector.Row(userButton))
	}

	selector.Inline(
		append(rows, selector.Row(btnPrev, btnNext))...,
	)

	b.Handle(&btnPrev, func(c *tb.Callback) {
		usersCount, _ := u.GetUsersCount()
		pagesCount := usersCount/userPerPage + 1

		page, _ := strconv.Atoi(c.Data)

		page--

		if page < 1 {
			page = int(pagesCount)
		}

		send := c.Message

		selector := getAdminUserSelector(page, m, b)

		_, err := b.Edit(send, getAdminUserMessage(page), selector, tb.ModeHTML)
		if err != nil {
			bot.SendMessageLog(err.Error(), b)
		}

		err = b.Respond(c, &tb.CallbackResponse{})
		if err != nil {
			bot.SendMessageLog(err.Error(), b)
		}
	})

	b.Handle(&btnNext, func(c *tb.Callback) {
		usersCount, _ := u.GetUsersCount()
		pagesCount := usersCount/userPerPage + 1

		page, _ := strconv.Atoi(c.Data)

		page++

		if page > int(pagesCount) {
			page = 1
		}

		send := c.Message

		selector := getAdminUserSelector(page, m, b)

		_, err := b.Edit(send, getAdminUserMessage(page), selector, tb.ModeHTML)
		if err != nil {
			bot.SendMessageLog(err.Error(), b)
		}

		err = b.Respond(c, &tb.CallbackResponse{})
		if err != nil {
			bot.SendMessageLog(err.Error(), b)
		}
	})

	return selector
}
