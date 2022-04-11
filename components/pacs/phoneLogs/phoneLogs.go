package phoneLogs

import (
	"encoding/json"
	"errors"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"os"
	"skandigatebot/base"
	"skandigatebot/bot"
	pc "skandigatebot/components/pacs/config"
	"skandigatebot/models/account"
	"skandigatebot/models/gateLog"
	"skandigatebot/models/gateLog/result"
	gateLogType "skandigatebot/models/gateLog/type"
	"skandigatebot/models/user"
	"strconv"
	"time"
)

func UpdateLogs(b *tb.Bot) {
	log.Print("Run update logs")

	maxNumber := gateLog.GetLastPhoneLogNumber()

	pacsLogs, err := getPACSLogs(maxNumber)

	if err != nil {
		fmt.Println(err)

		return
	}

	addLogs(pacsLogs, b)
	UpdateNonPhoneLogsOpenAt()

	log.Print("Finish update logs")
}

func getPACSLogs(fromNumber uint) (PACSLogResponse, error) {
	conf := pc.New()

	client := &http.Client{}
	URL := conf.Host + "/data.cgx?cmd={\"Command\":\"GetNextGateRecords\",\"Index\":" + strconv.FormatUint(uint64(fromNumber), 10) + "}"

	req, err := http.NewRequest("GET", URL, nil)
	req.SetBasicAuth(conf.User, conf.Password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

		return PACSLogResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp)

		return PACSLogResponse{}, errors.New("Кривой статус код")
	}

	pacUsersResponse := &PACSLogResponse{}
	err = json.NewDecoder(resp.Body).Decode(&pacUsersResponse)
	if err != nil {
		fmt.Println(err)

		return PACSLogResponse{}, err
	}

	return *pacUsersResponse, err
}

func addLogs(pacsLogs PACSLogResponse, b *tb.Bot) {
	for _, pacsLog := range pacsLogs.Records {
		logNumberInt, err := strconv.Atoi(fmt.Sprintf("%v", pacsLog[0]))

		if err != nil {
			continue
		}

		logNumber := uint(logNumberInt)

		logTime, err := time.Parse("2006-01-02T15:04:05", fmt.Sprintf("%v", pacsLog[1]))

		if err != nil {
			continue
		}

		logTime = logTime.Add(-3 * time.Hour)

		logPhoneInt, err := strconv.Atoi(fmt.Sprintf("%v", pacsLog[2]))

		if err != nil {
			continue
		}

		logPhone := uint(logPhoneInt)

		newGateLog := gateLog.GateLog{
			Phone:         logPhone,
			ResultId:      result.Success,
			LogTypeId:     gateLogType.Phone,
			LogTypeNumber: logNumber,
			OpenAt:        logTime,
		}

		foundUser, err := user.GetUser(logPhone)

		if err != user.ErrNotFound {
			newGateLog.UserId = foundUser.Id
		}

		err = base.GetDB().Create(&newGateLog).Error

		if err != nil {
			fmt.Println(err)
		}

		foundAccount, err := account.GetAccountByPhone(logPhone)

		logMessage := ""
		logMessage += os.Getenv("ENV")
		logMessage += " :: " + (logTime).Format("2006-01-02 15:04:05")

		if logPhone != 0 {
			logMessage += " :: +" + strconv.Itoa(int(logPhone))
		}

		logMessage += "\n"
		logMessage += "☎️ "
		logMessage += "<a href=\"tg://user?id=" + strconv.FormatInt(int64(foundAccount.AccountId), 10) + "\">"

		logMessage += foundAccount.FirstName
		logMessage += " "
		logMessage += foundAccount.LastName

		if foundAccount.UserName != "" {
			logMessage += " ("
			logMessage += foundAccount.UserName
			logMessage += ")"
		}

		logMessage += "</a> "

		logMessage += "open gate"
		logMessage = "✅ " + logMessage

		bot.SendMessageLog(logMessage, b)
	}
}

func UpdateNonPhoneLogsOpenAt() {
	base.GetDB().Exec("UPDATE tg_gate_log SET open_at = created_at WHERE phone is NULL and open_at <> created_at;")
}
