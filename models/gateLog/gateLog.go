package gateLog

import (
	"errors"
	"gorm.io/gorm"
	"skandigatebot/base"
	"skandigatebot/models"
	"skandigatebot/models/gateLog/result"
	gateLogType "skandigatebot/models/gateLog/type"
	"skandigatebot/models/orm"
	"skandigatebot/models/user"
	"time"
)

type GateLog struct {
	orm.AutoId
	UserId        uint                    `gorm:"default:0" json:"-"`
	User          user.User               `json:"-"`
	ResultId      uint                    `gorm:"default:1" json:"-"`
	Result        result.Result           `json:"-"`
	LogTypeId     uint                    `gorm:"default:1" json:"-"`
	LogType       gateLogType.GateLogType `json:"-"`
	LogTypeNumber uint                    `gorm:"default:0" json:"-"`
	Phone         uint
	OpenAt        time.Time `gorm:"type:datetime;default:now();index" json:"-"`
	orm.Time
}

var (
	ErrNotFound = errors.New("log not found")
)

func (gl *GateLog) BeforeCreate(tx *gorm.DB) (err error) {
	gl.CreatedAt = time.Now()
	gl.UpdatedAt = time.Now()

	return
}

func (gl *GateLog) BeforeUpdate(tx *gorm.DB) (err error) {
	gl.UpdatedAt = time.Now()

	return
}

func LogSuccess(userId uint) {
	gl := &GateLog{}

	gl.UserId = userId
	gl.ResultId = result.Success
	gl.OpenAt = time.Now()

	base.GetDB().Create(&gl)
}

func LogFail(userId uint) {
	gl := &GateLog{}

	gl.UserId = userId
	gl.ResultId = result.Fail
	gl.OpenAt = time.Now()

	base.GetDB().Create(&gl)
}

func GetLogsCount() (int64, error) {
	var logsCount int64

	base.
		GetDB().
		Model(&GateLog{}).
		Count(&logsCount)

	return logsCount, nil
}

func GetLogsWithUsers(offset int, limit int) ([]models.LogUserAccount, error) {
	var logs []models.LogUserAccount

	res := base.
		GetDB().
		Model(&GateLog{}).
		Select(
			"tg_gate_log.id as LogId," +
				"tg_gate_log.open_at as LogOpenAt," +
				"tg_gate_log.result_id as LogResultId," +
				"tg_gate_log.log_type_id as LogTypeId," +
				"tg_gate_log.log_type_number as LogTypeNumber," +
				"tg_user.id as UserId," +
				"tg_user.first_name as UserFirstName," +
				"tg_user.phone as phone," +
				"tg_user.role_id as RoleId," +
				"tg_user.active_id as ActiveId," +
				"tg_account.first_name as AccountFirstName," +
				"tg_account.last_name as AccountLastName," +
				"tg_account.user_name as AccountUserName").
		Joins("left join tg_user on tg_user.id = tg_gate_log.user_id").
		Joins("left join tg_account on tg_account.phone = tg_user.phone").
		Offset(offset).
		Limit(limit).
		Order("tg_gate_log.open_at desc").
		Find(&logs)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return logs, ErrNotFound
	}

	return logs, nil
}

func GetLastPhoneLogNumber() uint {
	var log GateLog

	res := base.
		GetDB().
		Model(&GateLog{}).
		Order("log_type_number DESC").
		Limit(1).
		First(&log)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0
	}

	return log.LogTypeNumber
}
