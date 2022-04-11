package account

import (
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
	"skandigatebot/base"
	"skandigatebot/models/orm"
	"time"
)

type Account struct {
	orm.AutoId
	AccountId uint   `gorm:"uniqueIndex"`
	FirstName string `gorm:"type:varchar(100)"`
	LastName  string `gorm:"type:varchar(100)"`
	UserName  string `gorm:"type:varchar(100)"`
	Phone     uint   `gorm:"index"`
	orm.Time
}

var (
	ErrNotFound = errors.New("account not found")
)

func (account *Account) BeforeCreate(tx *gorm.DB) (err error) {
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	return
}

func (account *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	account.UpdatedAt = time.Now()

	return
}

func GetAccountByPhone(phone uint) (Account, error) {
	var account Account

	result := base.
		GetDB().
		Model(&Account{}).
		Where("phone = ?", phone).
		Take(&account)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return account, ErrNotFound
	}

	return account, nil
}

func GetAccount(m *tb.Message) Account {
	accountId := m.Sender.ID

	var account Account

	result := base.
		GetDB().
		Model(&Account{}).
		Where("account_id = ?", accountId).
		Take(&account)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		account = Account{
			AccountId: uint(accountId),
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			UserName:  m.Sender.Username,
		}

		base.GetDB().Save(&account)
	} else {
		hasChanges := false
		if account.FirstName != m.Sender.FirstName {
			account.FirstName = m.Sender.FirstName
			hasChanges = true
		}
		if account.LastName != m.Sender.LastName {
			account.LastName = m.Sender.LastName
			hasChanges = true
		}
		if account.UserName != m.Sender.Username {
			account.UserName = m.Sender.Username
			hasChanges = true
		}

		if hasChanges {
			base.GetDB().Save(&account)
		}
	}

	return account
}
