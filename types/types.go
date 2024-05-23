package types

import "gorm.io/gorm"

type BasketState int32

// * gorm can't update zero values :_)
const (
	BasketState_PENDING BasketState = iota + 1
	BasketState_COMPLETED
)

// type BasketStateValue struct {
// 	Value BasketState
// }

// type Int32Value struct {
// 	Value int32
// }

// func (x *Int32Value) GetValue() int32 {
// 	if x != nil {
// 		return x.Value
// 	}
// 	return 0
// }

type Basket struct {
	gorm.Model
	Data  string `gorm:"size:2048"`
	State BasketState
}
