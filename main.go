package main

import (
	"errors"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
	"skandigatebot/base"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
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
	loadEnv()
	initSettings()

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_APITOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)

		return
	}

	b.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		}

		_, err := b.Send(m.Sender, textHello)
		if err != nil {
			log.Fatal(err)
		}

		user := getUser(m)

		checkAuth(user, m, b)
	})

	b.Handle(tb.OnContact, func(m *tb.Message) {
		user := getUser(m)

		if user.UserId != int(m.Contact.UserID) {
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

		if user.Phone != phone {
			user.Phone = phone

			base.GetDB().Save(&user)
		}

		checkAuth(user, m, b)
	})

	b.Start()
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func initSettings() {

}

func getUser(m *tb.Message) u.User {
	userId := m.Sender.ID

	var user u.User

	result := base.
		GetDB().
		Model(&u.User{}).
		Where("user_id = ?", userId).
		Take(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user = u.User{
			UserId:    int(userId),
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			UserName:  m.Sender.Username,
		}

		base.GetDB().Save(&user)
	} else {
		hasChanges := false
		if user.FirstName != m.Sender.FirstName {
			user.FirstName = m.Sender.FirstName
			hasChanges = true
		}
		if user.LastName != m.Sender.LastName {
			user.LastName = m.Sender.LastName
			hasChanges = true
		}
		if user.UserName != m.Sender.Username {
			user.UserName = m.Sender.Username
			hasChanges = true
		}

		if hasChanges {
			base.GetDB().Save(&user)
		}
	}

	return user
}

func showAuthMenu(m *tb.Message, b *tb.Bot) {
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnSharePhoneNumber := menu.Contact(textSharePhoneNumber)

	menu.Reply(
		menu.Row(btnSharePhoneNumber),
	)

	_, err := b.Send(m.Sender, textNeedAuth, menu)
	if err != nil {
		log.Fatal(err)
	}
}

func showGateMenu(m *tb.Message, b *tb.Bot) {
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnOpenGate := menu.Text(textOpenGate)

	menu.Reply(
		menu.Row(btnOpenGate),
	)

	_, err := b.Send(m.Sender, textAlreadyAuth, menu)
	if err != nil {
		log.Fatal(err)
	}
}

func checkAuth(user u.User, m *tb.Message, b *tb.Bot) {
	if user.Phone > 0 {
		showGateMenu(m, b)
	} else {
		showAuthMenu(m, b)
	}
}
