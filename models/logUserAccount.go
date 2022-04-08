package models

import "time"

type LogUserAccount struct {
	LogId            uint      `gorm:"type:integer"`
	LogOpenAt        time.Time `gorm:"type:datetime;default:now()" json:"-"`
	LogResultId      uint      `gorm:"type:integer"`
	LogTypeId        uint      `gorm:"type:integer"`
	LogTypeNumber    uint      `gorm:"type:integer"`
	UserId           uint      `gorm:"type:integer"`
	UserFirstName    string    `gorm:"type:varchar(100)"`
	Phone            uint      `gorm:"uniqueIndex"`
	RoleId           uint      `gorm:"default:2" json:"-"`
	ActiveId         uint      `gorm:"default:1" json:"-"`
	AccountFirstName string    `gorm:"type:varchar(100)"`
	AccountLastName  string    `gorm:"type:varchar(100)"`
	AccountUserName  string    `gorm:"type:varchar(100)"`
}
