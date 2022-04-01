package gateLogType

import (
	"gorm.io/gorm"
	"skandigatebot/base"
	"skandigatebot/models/orm"
	"time"
)

type GateLogType struct {
	orm.AutoId
	Name string `gorm:"uniqueIndex;type:varchar(100)"`
	orm.Time
}

const (
	Bot   = 1
	Phone = 2
)

func (glt *GateLogType) BeforeCreate(tx *gorm.DB) (err error) {
	glt.CreatedAt = time.Now()
	glt.UpdatedAt = time.Now()

	return
}

func (r *GateLogType) BeforeUpdate(tx *gorm.DB) (err error) {
	r.UpdatedAt = time.Now()

	return
}

func SeedGateLogTypes() {
	base.GetDB().Create(&GateLogType{
		Name: "Бот",
	})

	base.GetDB().Create(&GateLogType{
		Name: "Звонок",
	})
}
