package main

import (
	"log"
	"os"
	"skandigatebot/console"
	"skandigatebot/screens/admin"
	"skandigatebot/screens/auth"
	"skandigatebot/screens/first"
	"skandigatebot/screens/gate"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
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

	b.Start()
}
