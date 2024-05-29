package main

import (
	"fmt"
	"net/http"
	db_manager "shop-basket/db"
	"shop-basket/types"
	"shop-basket/utils"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
	_ "gorm.io/gorm"
)

func main() {
	if godotenv.Load() != nil {
		fmt.Println("load env file error")
		return
	}

	e := echo.New()
	db := db_manager.InitDB()

	fmt.Println("start service")

	e.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello World")
	})

	e.GET("basket", func(c echo.Context) error {
		var baskets []types.Basket
		if res := db.Find(&baskets); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}

		return c.JSON(http.StatusOK, baskets)
	})

	e.GET("basket/:id", func(c echo.Context) error {
		id, err := utils.GetIntParam(c, "id")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		var basket types.Basket
		if res := db.First(&basket, id); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}
		return c.JSON(http.StatusOK, basket)
	})

	e.POST("basket", func(c echo.Context) error {
		//todo: use validator
		data := c.FormValue("data")
		if len(data) > 2048 {
			return c.JSON(http.StatusBadRequest, "data length limit exceeded!!")
		}

		stateValue, err, isEmpty := utils.GetIntValue(c, "state")
		if !isEmpty && err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		state, ok := types.BasketStateMap[int32(stateValue)]
		if isEmpty || !ok {
			state = types.BasketState_PENDING
		}

		basket := types.Basket{
			Data:  datatypes.JSON(data),
			State: state,
			Items: types.Items{},
		}

		if res := db.Create(&basket); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}
		return c.JSON(http.StatusOK, basket.ID)
	})

	e.PATCH("basket/:id", func(c echo.Context) error {
		id, err := utils.GetIntParam(c, "id")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		data, isEmpty := utils.GetStringValue(c, "data")
		modifyData := !isEmpty

		stateValue, err, isEmpty := utils.GetIntValue(c, "state")
		modifyState := !isEmpty
		if !isEmpty && err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		name, isNameEmpty := utils.GetStringValue(c, "name")
		price, err, isPriceEmpty := utils.GetIntValue(c, "price")
		if !isEmpty && err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		shouldAddItem := !isNameEmpty && !isPriceEmpty

		//todo: enum range check with map

		//todo: use proto buff like values
		state := types.BasketState(stateValue)

		var oldBasket types.Basket
		if res := db.First(&oldBasket, id); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}

		if oldBasket.State == types.BasketState_COMPLETED {
			return c.JSON(http.StatusLocked, "completed basket can not be changed")
		}

		var newBasket types.Basket
		if modifyData {
			newBasket.Data = datatypes.JSON(data)
		}
		if modifyState {
			newBasket.State = state
		}
		if shouldAddItem {
			newBasket.Items = append(oldBasket.Items, types.Item{Name: name, Price: price})
		}

		if res := db.Model(&types.Basket{}).Where("ID = ?", id).Updates(&newBasket); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}

		return c.JSON(http.StatusOK, "Basket successfully modified!")
	})

	e.DELETE("basket/:id", func(c echo.Context) error {
		id, err := utils.GetIntParam(c, "id")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		//* hard delete
		if res := db.Unscoped().Delete(&types.Basket{}, id); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}
		return c.JSON(http.StatusOK, "Basket deleted successfully!")
	})

	e.Logger.Fatal(e.Start(":3005"))
}
