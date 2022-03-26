package active

import (
	"gorm.io/gorm"
	"skandigatebot/base"
	"skandigatebot/models/orm"
	"time"
)

type Active struct {
	orm.AutoId
	Name string `gorm:"uniqueIndex;type:varchar(100)"`
	orm.Time
}

const (
	Allow   = 1
	Blocked = 2
)

func (active *Active) BeforeCreate(tx *gorm.DB) (err error) {
	active.CreatedAt = time.Now()
	active.UpdatedAt = time.Now()

	return
}

func (active *Active) BeforeUpdate(tx *gorm.DB) (err error) {
	active.UpdatedAt = time.Now()

	return
}

func SeedActives() {
	base.GetDB().Create(&Active{
		Name: "Активный пользователь",
	})

	base.GetDB().Create(&Active{
		Name: "Заблокированный пользователь",
	})
}
