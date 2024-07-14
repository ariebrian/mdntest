package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

func InitDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=myappuser dbname=myapp password=yourpassword sslmode=disable")
	if err != nil {
		panic("failed to connect database"  + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}
