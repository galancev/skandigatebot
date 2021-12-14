package orm

import (
	"github.com/gofrs/uuid"
)

type Uuid struct {
	Id uuid.UUID `gorm:"primary_key;type:char(36);not null" json:"id"`
}
