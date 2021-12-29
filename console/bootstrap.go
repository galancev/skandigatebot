package console

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"skandigatebot/base"
	a "skandigatebot/models/account"
	"skandigatebot/models/user"
	"skandigatebot/models/user/role"
)

func Boot() {
	loadEnv()
	initSettings()
	handleArgs()
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func initSettings() {
	_ = base.GetDB().Debug().AutoMigrate(&a.Account{}, &user.User{}, &role.Role{})
	role.SeedRoles()
	user.SeedUsers()
}

func handleArgs() {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {

		case "seed":
			seed(args)

		}

		os.Exit(0)
	}
}

func seed(args []string) {
	if len(args) >= 2 {
		switch args[1] {

		case "roles":
			role.SeedRoles()

		case "users":
			user.SeedUsers()

		}
	} else {
		role.SeedRoles()
		user.SeedUsers()
	}
}
