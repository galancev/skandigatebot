package phoneLogs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"skandigatebot/base"
	pc "skandigatebot/components/pacs/config"
	"skandigatebot/models/gateLog"
	"skandigatebot/models/gateLog/result"
	gateLogType "skandigatebot/models/gateLog/type"
	"strconv"
	"time"
)

func UpdateLogs() {
	log.Print("Run update logs")

	maxNumber := gateLog.GetLastPhoneLogNumber()

	pacsLogs, err := getPACSLogs(maxNumber)

	if err != nil {
		fmt.Println(err)

		return
	}

	addLogs(pacsLogs)
	UpdateUserIds()

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

func addLogs(pacsLogs PACSLogResponse) {
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

		err = base.GetDB().Create(&newGateLog).Error

		fmt.Println(err)
	}
}

func UpdateUserIds() {
	base.GetDB().Exec("UPDATE tg_gate_log gl INNER JOIN tg_user tu ON gl.phone = tu.phone SET gl.user_id = tu.id WHERE gl.log_type_id = 2 AND gl.user_id = 0;")
}
