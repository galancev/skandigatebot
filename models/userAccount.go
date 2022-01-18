package models

type UserAccount struct {
	UserId           uint   `gorm:"type:integer"`
	UserFirstName    string `gorm:"type:varchar(100)"`
	UserLastName     string `gorm:"type:varchar(100)"`
	Phone            uint   `gorm:"uniqueIndex"`
	RoleId           uint   `gorm:"default:1" json:"-"`
	AccountFirstName string `gorm:"type:varchar(100)"`
	AccountLastName  string `gorm:"type:varchar(100)"`
	AccountUserName  string `gorm:"type:varchar(100)"`
}
