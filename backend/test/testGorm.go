package main

import (
	"ginchat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3307)/ginchat?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.GroupBasic{})
	db.AutoMigrate(&models.Contact{})
	db.AutoMigrate(&models.Community{})

	// user := &models.UserBasic{}
	// user.Name = "測試a"
	// db.Create(user)

	// fmt.Println(db.First(user, 1))

	// db.Model(user).Update("PassWord", "1234")
}
