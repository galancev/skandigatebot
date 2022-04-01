package main

import (
	"context"
	"github.com/procyon-projects/chrono"
	"log"
	"os"
	"skandigatebot/bot"
	"skandigatebot/components/pacs/users"
	"skandigatebot/console"
	"skandigatebot/screens/admin"
	"skandigatebot/screens/auth"
	"skandigatebot/screens/first"
	"skandigatebot/screens/gate"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	al "skandigatebot/screens/admin/log"
	au "skandigatebot/screens/admin/users"
)

const (
	textUnknownText     = "😅 Не понял вас."
	textUnknownPhoto    = "😱 Картинки, фоточки, ура! Сохраню в свой альбом, буду разглядывать на досуге. 👌"
	textUnknownVideo    = "❤️ Спасибо! Мне было так скучно..."
	textUnknownDocument = "😏 Что мне с этим делать?"
)

func main() {
	/*defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()*/

	console.Boot()

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_APITOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)

		return
	}

	scheduler()

	bot.SendMessageLog("Bot starting...", b)

	b.Handle("/start", func(m *tb.Message) {
		pauth := auth.New()
		pgate := gate.New(pauth)
		pfirst := first.New(pauth, pgate)

		pfirst.OnStart(m, b)
	})

	b.Handle(tb.OnContact, func(m *tb.Message) {
		pauth := auth.New()

		pauth.OnAuth(m, b)
	})

	b.Handle(gate.OpenGateButton, func(m *tb.Message) {
		pauth := auth.New()
		pgate := gate.New(pauth)

		pgate.OnOpen(m, b)
	})

	b.Handle(admin.OnAdminButton, func(m *tb.Message) {
		pauth := auth.New()
		pgate := gate.New(pauth)
		padmin := admin.New(pauth, pgate)

		padmin.OnAdmin(m, b)
	})

	b.Handle(admin.OnAdminExitButton, func(m *tb.Message) {
		account, user, _ := bot.GetAccountAndUser(m)

		pauth := auth.New()
		pgate := gate.New(pauth)
		pfirst := first.New(pauth, pgate)

		pfirst.ShowFirstMenu(&account, &user, m, b)
	})

	b.Handle(admin.OnAdminShowUsers, func(m *tb.Message) {
		pauth := auth.New()
		pgate := gate.New(pauth)
		pau := au.New(pauth, pgate)

		pau.OnAdminUsers(m, b)
	})

	b.Handle(admin.OnAdminShowLog, func(m *tb.Message) {
		pauth := auth.New()
		pgate := gate.New(pauth)
		pal := al.New(pauth, pgate)

		pal.OnAdminLog(m, b)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		bot.SendMessage(textUnknownText, m, b)

		account, user, _ := bot.GetAccountAndUser(m)

		pauth := auth.New()
		pgate := gate.New(pauth)
		pfirst := first.New(pauth, pgate)

		pfirst.ShowFirstMenu(&account, &user, m, b)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		bot.SendMessage(textUnknownPhoto, m, b)

		account, user, _ := bot.GetAccountAndUser(m)

		pauth := auth.New()
		pgate := gate.New(pauth)
		pfirst := first.New(pauth, pgate)

		pfirst.ShowFirstMenu(&account, &user, m, b)
	})

	b.Handle(tb.OnVideo, func(m *tb.Message) {
		bot.SendMessage(textUnknownVideo, m, b)

		account, user, _ := bot.GetAccountAndUser(m)

		pauth := auth.New()
		pgate := gate.New(pauth)
		pfirst := first.New(pauth, pgate)

		pfirst.ShowFirstMenu(&account, &user, m, b)
	})

	b.Handle(tb.OnDocument, func(m *tb.Message) {
		bot.SendMessage(textUnknownDocument, m, b)

		account, user, _ := bot.GetAccountAndUser(m)

		pauth := auth.New()
		pgate := gate.New(pauth)
		pfirst := first.New(pauth, pgate)

		pfirst.ShowFirstMenu(&account, &user, m, b)
	})

	b.Handle(tb.OnQuery, func(m *tb.Message) {
		log.Print(m.Text)
		// all the text messages that weren't
		// captured by existing handlers
	})

	b.Handle(tb.OnCallback, func(m *tb.Message) {
		log.Print(m.Text)
		// all the text messages that weren't
		// captured by existing handlers
	})

	b.Start()
}

func scheduler() {
	var err error

	taskScheduler := chrono.NewDefaultTaskScheduler()

	_, err = taskScheduler.ScheduleWithCron(func(ctx context.Context) {
		go users.UpdateUsers()
	}, "0 0 * * * *")

	if err == nil {
		log.Print("Task has been scheduled")
	}

	go users.UpdateUsers()
	//go phoneLogs.UpdateLogs()
}
