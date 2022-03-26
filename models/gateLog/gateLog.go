package gateLog

import (
	"gorm.io/gorm"
	"skandigatebot/base"
	"skandigatebot/models/gateLog/result"
	"skandigatebot/models/orm"
	"skandigatebot/models/user"
	"time"
)

type GateLog struct {
	orm.AutoId
	UserId   uint          `gorm:"default:0" json:"-"`
	User     user.User     `json:"-"`
	ResultId uint          `gorm:"default:1" json:"-"`
	Result   result.Result `json:"-"`
	orm.Time
}

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

	base.GetDB().Create(&gl)
}

func LogFail(userId uint) {
	gl := &GateLog{}

	gl.UserId = userId
	gl.ResultId = result.Fail

	base.GetDB().Create(&gl)
}
