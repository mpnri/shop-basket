package main

import (
	"fmt"
	"net/http"
	db_manager "shop-basket/db"
	"shop-basket/types"
	"shop-basket/utils"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "gorm.io/gorm"
)

func main() {
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
		//todo: validate data
		data := c.FormValue("data")
		// state, err := strconv.Atoi(c.FormValue("state"))
		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, err)
		// }
		basket := types.Basket{Data: data, State: types.BasketState_PENDING}

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

		data := c.FormValue("data")
		modifyData := true
		if data == "" {
			modifyData = false
		}

		stateValue := c.FormValue("state")
		intState, err := strconv.Atoi(stateValue)
		if stateValue != "" && err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		//todo: enum range check with map

		//todo: use proto buff like values
		state := types.BasketState(intState)
		modifyState := true
		if stateValue == "" {
			modifyState = false
		}

		var oldBasket types.Basket
		if res := db.First(&oldBasket, id); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}

		if oldBasket.State == types.BasketState_COMPLETED {
			return c.JSON(http.StatusLocked, "completed basket can not be changed")
		}

		var newBasket types.Basket
		if modifyData {
			newBasket.Data = data
		}
		if modifyState {
			newBasket.State = state
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

		if res := db.Delete(&types.Basket{}, id); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}
		return c.JSON(http.StatusOK, "Basket deleted successfully!")
	})

	e.Logger.Fatal(e.Start(":3005"))
}
