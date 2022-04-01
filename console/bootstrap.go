package console

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"skandigatebot/base"
	a "skandigatebot/models/account"
	"skandigatebot/models/gateLog"
	"skandigatebot/models/gateLog/result"
	gateLogType "skandigatebot/models/gateLog/type"
	"skandigatebot/models/user"
	"skandigatebot/models/user/active"
	"skandigatebot/models/user/role"
)

func Boot() {
	loadEnv()
	initSettings()
	handleArgs()
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func initSettings() {
	_ = base.GetDB().Debug().AutoMigrate(&a.Account{}, &user.User{}, &role.Role{}, &active.Active{}, &result.Result{}, &gateLog.GateLog{})
	seed([]string{})
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
		active.SeedActives()
		user.SeedUsers()
		result.SeedGateResults()
		gateLogType.SeedGateLogTypes()

		base.GetDB().Exec("alter table tg_gate_log drop foreign key fk_tg_gate_log_user")
		base.GetDB().Exec("update tg_gate_log set open_at = created_at")
	}
}
