package db_manager

import (
	types "shop-basket/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Basket struct {
	gorm.Model
	Data  string `gorm:"size:2048"`
	State types.BasketState
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:data.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Basket{})

	return db
}
