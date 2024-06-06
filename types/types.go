package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	// "github.com/lib/pq"
	"github.com/golang-jwt/jwt/v5"
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

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
	// Baskets  pq.Int32Array `gorm:"type:integer[]"`
}

type Basket struct {
	gorm.Model
	UserID uint
	Data   datatypes.JSON `gorm:"size:2048"`
	Items  Items          `json:"items" gorm:"type:jsonb"`
	State  BasketState
}

type Item struct {
	ShoppingCartID uint
	Name           string `json:"name"`
	Price          int32  `json:"price"`
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

type JwtCustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}
