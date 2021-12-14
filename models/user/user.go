package user

import (
	"gorm.io/gorm"
	"skandigatebot/models/orm"
	"time"
)

type User struct {
	orm.AutoId
	UserId    int    `gorm:"uniqueIndex"`
	FirstName string `gorm:"type:varchar(100)"`
	LastName  string `gorm:"type:varchar(100)"`
	UserName  string `gorm:"type:varchar(100)"`
	Phone     int    `gorm:"uniqueIndex"`
	orm.Time
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now()

	return
}
