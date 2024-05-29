package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BasketState int32

// * gorm can't update zero values :_)
const (
	BasketState_PENDING BasketState = iota + 1
	BasketState_COMPLETED
)

var (
	BasketStateMap = map[int32]BasketState{
		1: BasketState_PENDING,
		2: BasketState_COMPLETED,
	}
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

type Item struct {
	Name  string `json:"name"`
	Price int32  `json:"price"`
}

type Items []Item

func (j *Items) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("failed to unmarshal JSONB value")
    }

    var items Items
    if err := json.Unmarshal(bytes, &items); err != nil {
        return err
    }

    *j = items
    return nil
}

func (j Items) Value() (driver.Value, error) {
    if len(j) == 0 {
        return nil, nil
    }
    return json.Marshal(j)
}


type Basket struct {
	gorm.Model
	Data  datatypes.JSON `gorm:"size:2048"`
	Items Items         `json:"items" gorm:"type:jsonb"`
	State BasketState
}
