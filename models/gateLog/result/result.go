package result

import (
	"gorm.io/gorm"
	"skandigatebot/base"
	"skandigatebot/models/orm"
	"time"
)

type Result struct {
	orm.AutoId
	Name string `gorm:"uniqueIndex;type:varchar(100)"`
	orm.Time
}

const (
	Success = 1
	Fail    = 2
)

func (r *Result) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()

	return
}

func (r *Result) BeforeUpdate(tx *gorm.DB) (err error) {
	r.UpdatedAt = time.Now()

	return
}

func SeedGateResults() {
	base.GetDB().Create(&Result{
		Name: "Успешное открытие",
	})

	base.GetDB().Create(&Result{
		Name: "Ошибка",
	})
}
