package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetStringParam(c echo.Context, name string) (string, bool) {
	value := c.Param(name)
	isEmpty := value == ""
	return value, isEmpty
}

func GetIntParam(c echo.Context, name string) (int, error) {
	value, isEmpty := GetStringParam(c, name)
	if isEmpty {
		return 0, nil
	}

	res, err := strconv.Atoi(value)
	return res, err
}

func GetStringValue(c echo.Context, name string) (string, bool) {
	value := c.FormValue(name)
	isEmpty := value == ""
	return value, isEmpty
}

func GetIntValue(c echo.Context, name string) (int, error, bool) {
	value, isEmpty := GetStringValue(c, name)
	if isEmpty {
		return 0, nil, isEmpty
	}

	res, err := strconv.Atoi(value)
	return res, err, isEmpty
}

