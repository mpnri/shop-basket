package main

import (
	"fmt"
	"net/http"
	"os"
	db_manager "shop-basket/db"
	"shop-basket/types"
	"shop-basket/utils"
	_ "strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello World")
	},
	)

	e.POST("/register", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		user := types.User{Username: username, Password: password}
		db.Create(&user)

		return c.JSON(http.StatusCreated, user.ID)
	})
	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		var user types.User
		if err := db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
			return echo.ErrUnauthorized
		}

		claims := types.JwtCustomClaims{
			UserID: user.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(os.Getenv("SECRET")))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	})

	r := e.Group("/")
	_ = echojwt.ErrJWTMissing
	r.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(os.Getenv("SECRET")),
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return &types.JwtCustomClaims{} },
	}))

	r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "JWT token missing or invalid")
			}

			claims, ok := token.Claims.(*types.JwtCustomClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "failed to cast claims as JwtCustomClaims")
			}

			uid := claims.UserID
			var user types.User

			if res := db.Find(&user, uid); res.Error != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "user not found")
			}
			c.Set("uid", uid)
			return next(c)
		}
	})

	r.GET("basket", func(c echo.Context) error {
		uid, ok := c.Get("uid").(uint)
		if !ok {
			return echo.ErrUnauthorized
		}
		var baskets []types.Basket
		if res := db.Limit(10).Find(&types.Basket{UserID: uid}).Find(&baskets); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}

		return c.JSON(http.StatusOK, baskets)
	})

	r.GET("basket/:id", func(c echo.Context) error {
		user_id, ok := c.Get("uid").(uint)
		if !ok {
			return echo.ErrUnauthorized
		}

		id, err := utils.GetIntParam(c, "id")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		var basket types.Basket
		if res := db.First(&basket, id); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}
		if basket.UserID != user_id {
			return echo.ErrForbidden
		}

		return c.JSON(http.StatusOK, basket)
	})

	r.POST("basket", func(c echo.Context) error {
		user_id, ok := c.Get("uid").(uint)
		if !ok {
			return echo.ErrUnauthorized
		}

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
			Data:   datatypes.JSON(data),
			State:  state,
			Items:  types.Items{},
			UserID: user_id,
		}

		if res := db.Create(&basket); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}
		return c.JSON(http.StatusOK, basket.ID)
	})

	r.PATCH("basket/:id", func(c echo.Context) error {
		user_id, ok := c.Get("uid").(uint)
		if !ok {
			return echo.ErrUnauthorized
		}

		id, err := utils.GetIntParam(c, "id")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		var oldBasket types.Basket
		if res := db.First(&oldBasket, id); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}

		if oldBasket.UserID != user_id {
			return echo.ErrForbidden
		}

		if oldBasket.State == types.BasketState_COMPLETED {
			return c.JSON(http.StatusLocked, "completed basket can not be changed")
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

	r.DELETE("basket/:id", func(c echo.Context) error {
		user_id, ok := c.Get("uid").(uint)
		if !ok {
			return echo.ErrUnauthorized
		}

		id, err := utils.GetIntParam(c, "id")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		//* hard delete
		if res := db.Unscoped().Delete(&types.Basket{UserID: user_id}, id); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}
		return c.JSON(http.StatusOK, "Basket deleted successfully!")
	})

	e.Logger.Fatal(e.Start(":3005"))
}
