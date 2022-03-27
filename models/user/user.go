package user

import (
	"errors"
	"gorm.io/gorm"
	"skandigatebot/base"
	"skandigatebot/models"
	"skandigatebot/models/orm"
	"skandigatebot/models/user/active"
	"skandigatebot/models/user/role"
	"time"
)

type User struct {
	orm.AutoId
	FirstName string        `gorm:"type:varchar(100)"`
	Phone     uint          `gorm:"uniqueIndex"`
	RoleId    uint          `gorm:"default:2" json:"-"`
	Role      role.Role     `json:"-"`
	ActiveId  uint          `gorm:"default:1" json:"-"`
	Active    active.Active `json:"-"`
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
	return user.RoleId == role.Admin
}

func (user *User) IsUser() bool {
	return user.RoleId == role.User
}

func (user *User) IsActive() bool {
	return user.ActiveId == active.Allow
}

func (user *User) IsBlocked() bool {
	return user.ActiveId == active.Blocked
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
		RoleId:    role.Admin,
		ActiveId:  active.Allow,
	})

	base.GetDB().Create(&User{
		Phone:     79958848775,
		FirstName: "Алексей",
		RoleId:    role.Admin,
		ActiveId:  active.Allow,
	})
}

func GetUsersCount() (int64, error) {
	var userCount int64

	base.
		GetDB().
		Model(&User{}).
		Count(&userCount)

	return userCount, nil
}

func GetUsersWithAccount(offset int, limit int) ([]models.UserAccount, error) {
	var users []models.UserAccount

	result := base.
		GetDB().
		Model(&User{}).
		Select(
			"tg_user.id as UserId," +
				"tg_user.first_name as UserFirstName," +
				"tg_user.phone as phone," +
				"tg_user.role_id as RoleId," +
				"tg_user.active_id as ActiveId," +
				"tg_account.first_name as AccountFirstName," +
				"tg_account.last_name as AccountLastName," +
				"tg_account.user_name as AccountUserName").
		Joins("left join tg_account on tg_account.phone = tg_user.phone").
		Offset(offset).
		Limit(limit).
		Order("IF(tg_account.id IS NULL, 1, 0), tg_user.active_id, tg_user.id").
		Find(&users)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return users, ErrNotFound
	}

	return users, nil
}

func GetUsers() ([]User, error) {
	var users []User

	result := base.
		GetDB().
		Model(&User{}).
		Find(&users)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return users, ErrNotFound
	}

	return users, nil
}
