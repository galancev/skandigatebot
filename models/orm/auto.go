package orm

type AutoId struct {
	Id uint `gorm:"primarykey" json:"id"`
}
