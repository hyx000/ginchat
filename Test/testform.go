package main

import "gorm.io/gorm"

import (
	"ginchat/models"
	"gorm.io/driver/mysql"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(mysql.Open("root:root123@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		//panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&models.GroupBasic{})
	db.AutoMigrate(&models.Contact{})
}
