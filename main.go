package main

import (
	"log"
	"os"
	"skandigatebot/base"
	"skandigatebot/console"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	a "skandigatebot/models/account"
	u "skandigatebot/models/user"
)

const (
	textHello            = "Привет! Это скандибот для управления шлагбаумом паркинга 1 корпуса."
	textEnemyPhoneNumber = "Необходимо делиться своим телефоном, а не чужим!"
	textSharePhoneNumber = "Поделиться номером телефона"
	textNeedAuth         = "Для продолжения работы необходимо авторизоваться."
	textAlreadyAuth      = "Вы авторизованы, можете пользоваться шлагбаумом."
	textOpenGate         = "Открыть врата!"
)

func main() {
	console.Boot()

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_APITOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)

		return
	}

	b.Handle("/start", func(m *tb.Message) {
		_, err := b.Send(m.Sender, textHello)
		if err != nil {
			log.Fatal(err)
		}

		showFirstMenu(m, b)
	})

	b.Handle(tb.OnContact, func(m *tb.Message) {
		account := a.GetAccount(m)

		if account.AccountId != uint(m.Contact.UserID) {
			_, err := b.Send(m.Sender, textEnemyPhoneNumber)
			if err != nil {
				log.Fatal(err)
			}

			return
		}

		phone, err := strconv.Atoi(m.Contact.PhoneNumber)
		if err != nil {
			log.Fatal(err)
		}

		if account.Phone != uint(phone) {
			account.Phone = uint(phone)

			base.GetDB().Save(&account)
		}

		showFirstMenu(m, b)
	})

	b.Handle(textOpenGate, func(m *tb.Message) {
		account := a.GetAccount(m)

		if account.Phone > 0 {
			showGateMenu("Врата открываются!", m, b)
		} else {
			showAuthMenu("Вам нельзя это сделать, вы не авторизованы.", m, b)
		}
	})

	b.Start()
}

func showAuthMenu(message string, m *tb.Message, b *tb.Bot) {
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnSharePhoneNumber := menu.Contact(textSharePhoneNumber)

	menu.Reply(
		menu.Row(btnSharePhoneNumber),
	)

	_, err := b.Send(m.Sender, message, menu)
	if err != nil {
		log.Fatal(err)
	}
}

func showGateMenu(message string, m *tb.Message, b *tb.Bot) {
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnOpenGate := menu.Text(textOpenGate)

	menu.Reply(
		menu.Row(btnOpenGate),
	)

	_, err := b.Send(m.Sender, message, menu)
	if err != nil {
		log.Fatal(err)
	}
}

func showFirstMenu(m *tb.Message, b *tb.Bot) {
	account := a.GetAccount(m)

	if account.Phone > 0 {
		_, err := u.GetUser(account.Phone)

		if err != nil {
			if err == u.ErrNotFound {
				showAuthMenu("Вы успешно авторизовались, однако вашего телефона нет в списке разрешённых. Напишите скандифокс для добавления.", m, b)
			} else {
				showAuthMenu("Конина какая-то на сервере", m, b)
			}
		} else {
			showGateMenu(textAlreadyAuth, m, b)
		}
	} else {
		showAuthMenu(textNeedAuth, m, b)
	}
}
