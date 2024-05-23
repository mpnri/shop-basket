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

// var (
// 	BasketState_name = map[int32]string{
// 		0: "BasketState_COMPLETED",
// 		1: "BasketState_PENDING",
// 	}
// 	BasketState_value = map[string]int32{
// 		"BasketState_COMPLETED":    0,
// 		"BasketState_PENDING":     1,
// 	}
// )

type Basket struct {
	gorm.Model
	Data  string `gorm:"size:2048"`
	State BasketState
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Basket{})

	return db
}
