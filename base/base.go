package base

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"skandigatebot/config"
	u "skandigatebot/models/user"
	"strconv"
)

var db *gorm.DB

// Инициализация базы данных
func init() {
	conf := config.New()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Db.User, conf.Db.Password, conf.Db.Host, strconv.FormatInt(int64(conf.Db.Port), 10), conf.Db.Name)
	fmt.Println(dsn)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tg_",
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	//_ = db.Debug().Exec("set search_path = \"public\"")
	//_ = db.Debug().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\" SCHEMA public")
	_ = db.Debug().AutoMigrate(&u.User{})
}

// Возвращает объект бд
func GetDB() *gorm.DB {
	return db
}
