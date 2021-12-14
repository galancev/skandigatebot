package orm

import "time"

type Time struct {
	CreatedAt time.Time `gorm:"type:datetime;default:now()" json:"-"`
	UpdatedAt time.Time `gorm:"type:datetime;default:now()" json:"-"`
}
