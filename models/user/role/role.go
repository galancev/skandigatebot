package role

import (
	"gorm.io/gorm"
	"skandigatebot/base"
	"skandigatebot/models/orm"
	"time"
)

type Role struct {
	orm.AutoId
	Name string `gorm:"uniqueIndex;type:varchar(100)"`
	orm.Time
}

const (
	Admin = 1
	User  = 2
)

func (role *Role) BeforeCreate(tx *gorm.DB) (err error) {
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	return
}

func (role *Role) BeforeUpdate(tx *gorm.DB) (err error) {
	role.UpdatedAt = time.Now()

	return
}

func SeedRoles() {
	base.GetDB().Create(&Role{
		Name: "Администратор",
	})

	base.GetDB().Create(&Role{
		Name: "Пользователь",
	})
}
