package log

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"skandigatebot/bot"
	"skandigatebot/models"
	"skandigatebot/models/gateLog"
	"skandigatebot/models/gateLog/result"
	gateLogType "skandigatebot/models/gateLog/type"
	"strconv"
	"time"
	//au "skandigatebot/screens/admin/user"
)

const (
	textDbError          = "😵 Проблема с базой данных на сервере"
	textAuthAccessDenied = "❗️ Вы успешно авторизовались, однако вашего телефона нет в списке разрешенных. По вопросам доступа обратитесь в офис управляющей организации АО «ВК Комфорт» по адресу дом 1 корпус 3. Доступ возможен только для разгрузки и проезда на смежные территории."
	textAuthAdminDenied  = "📛 Хорошая попытка, но нет. В админку вам нельзя!"
	textNonAuth          = "⛔️ Вам нельзя это сделать, вы не авторизованы."

	userPerPage = 10
)

func getLogs(page int) ([]models.LogUserAccount, error) {
	logs, err := gateLog.GetLogsWithUsers((page-1)*userPerPage, userPerPage)

	return logs, err
}

func getAdminLogMessage(page int) string {
	usersCount, _ := gateLog.GetLogsCount()
	pagesCount := usersCount/userPerPage + 1

	var message string

	message += "Страница [" + strconv.Itoa(page) + "] из [" + strconv.Itoa(int(pagesCount)) + "]"

	return message
}

func getLogUserSelector(page int, m *tb.Message, b *tb.Bot) *tb.ReplyMarkup {
	selector := &tb.ReplyMarkup{}

	btnPrev := selector.Data("⬅", "prev", strconv.Itoa(page))
	btnNext := selector.Data("➡", "next", strconv.Itoa(page))

	users, err := getLogs(page)
	if err != nil {
		return nil
	}

	var userButtons []tb.Btn
	for _, user := range users {
		message := ""
		if user.LogResultId == result.Success {
			message += "✅"
		} else {
			message += "❌"
		}

		if user.LogTypeId == gateLogType.Bot {
			message += "🤖"
		} else {
			message += "☎️"
		}

		message += " " + (user.LogCreatedAt.Add(3 * time.Hour)).Format("2006-01-02 15:04:05")
		message += " +" + strconv.Itoa(int(user.Phone))
		message += " (" + user.UserFirstName + ")"
		if user.AccountFirstName != "" {
			message += " :: " + user.AccountFirstName + " " + user.AccountLastName
		}

		if user.AccountUserName != "" {
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
		usersCount, _ := gateLog.GetLogsCount()
		pagesCount := usersCount/userPerPage + 1

		page, _ := strconv.Atoi(c.Data)

		page--

		if page < 1 {
			page = int(pagesCount)
		}

		send := c.Message

		selector := getLogUserSelector(page, m, b)

		_, err := b.Edit(send, getAdminLogMessage(page), selector, tb.ModeHTML)
		if err != nil {
			bot.SendMessageLog(err.Error(), b)
		}

		err = b.Respond(c, &tb.CallbackResponse{})
		if err != nil {
			bot.SendMessageLog(err.Error(), b)
		}
	})

	b.Handle(&btnNext, func(c *tb.Callback) {
		usersCount, _ := gateLog.GetLogsCount()
		pagesCount := usersCount/userPerPage + 1

		page, _ := strconv.Atoi(c.Data)

		page++

		if page > int(pagesCount) {
			page = 1
		}

		send := c.Message

		selector := getLogUserSelector(page, m, b)

		_, err := b.Edit(send, getAdminLogMessage(page), selector, tb.ModeHTML)
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
