package user

import (
	"errors"
	"gorm.io/gorm"
	"skandigatebot/base"
	"skandigatebot/models/orm"
	"skandigatebot/models/user/role"
	"time"
)

type User struct {
	orm.AutoId
	FirstName string    `gorm:"type:varchar(100)"`
	LastName  string    `gorm:"type:varchar(100)"`
	Phone     uint      `gorm:"uniqueIndex"`
	RoleId    uint      `gorm:"default:1" json:"-"`
	Role      role.Role `json:"-"`
	orm.Time
}

var (
	ErrNotFound = errors.New("user not found")
)

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now()

	return
}

func (user *User) IsAdmin() bool {
	return user.Role.Id == role.Admin
}

func (user *User) IsUser() bool {
	return user.Role.Id == role.User
}

func GetUser(phone uint) (User, error) {
	var user User

	result := base.
		GetDB().
		Model(&User{}).
		Where("phone = ?", phone).
		Take(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, ErrNotFound
	}

	return user, nil
}

func SeedUsers() {
	base.GetDB().Create(&User{
		Phone:     79151019102,
		FirstName: "Евгений",
		LastName:  "Галанцев",
		RoleId:    role.Admin,
	})
}
