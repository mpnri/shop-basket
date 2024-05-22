package main

import (
	"net/http"
	db_manager "shop-basket/db"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db := db_manager.InitDB()

	e.GET("", func(context echo.Context) error {
		return context.JSON(http.StatusOK, "Hello World")
	})

	e.GET("basket", func(context echo.Context) error {
		var baskets []db_manager.Basket
		res := db.Find(&baskets)
		if res.Error != nil {
			return context.JSON(http.StatusInternalServerError, res.Error)
		}
		return context.JSON(http.StatusOK, baskets)
	})

	e.Logger.Fatal(e.Start(":3005"))
}
