package main

import (
	_ "fmt"
	"net/http"
	db_manager "shop-basket/db"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "gorm.io/gorm"
)

func main() {
	e := echo.New()
	db := db_manager.InitDB()

	e.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello World")
	})

	e.GET("basket", func(c echo.Context) error {
		var baskets []db_manager.Basket
		if res := db.Find(&baskets); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}

		return c.JSON(http.StatusOK, baskets)
	})

	e.GET("basket/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		var basket db_manager.Basket
		if res := db.First(&basket, id); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}
		return c.JSON(http.StatusOK, basket)
	})

	e.POST("basket", func(c echo.Context) error {
		data := c.FormValue("data")
		// state, err := strconv.Atoi(c.FormValue("state"))
		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, err)
		// }
		basket := db_manager.Basket{Data: data, State: db_manager.BasketState_PENDING}

		if res := db.Create(&basket); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}
		return c.JSON(http.StatusOK, basket.ID)
	})

	e.Logger.Fatal(e.Start(":3005"))
}
