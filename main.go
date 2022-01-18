package main

import (
	"log"
	"os"
	"skandigatebot/bot"
	"skandigatebot/console"
	"skandigatebot/screens/admin"
	"skandigatebot/screens/auth"
	"skandigatebot/screens/first"
	"skandigatebot/screens/gate"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	au "skandigatebot/screens/admin/users"
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

	b.Handle(tb.OnText, func(m *tb.Message) {
		log.Print(m.Text)
		// all the text messages that weren't
		// captured by existing handlers
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
