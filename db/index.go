package db_manager

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BasketState int32

const (
	BasketState_COMPLETED BasketState = 0
	BasketState_PENDING   BasketState = 1
)

type Basket struct {
	gorm.Model
	data  string `gorm:"size:2048"`
	state BasketState
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Basket{})

	return db
}
