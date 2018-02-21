package handler

import (
    "net/http"
    "github.com/labstack/echo"
)

func Hello() echo.HandlerFunc {
    return func(c echo.Context) error {     //c をいじって Request, Responseを色々する 
        return c.String(http.StatusOK, "Hello World")
    }
}
func Clac() echo.HandlerFunc {
    return func(c echo.Context) error {
        calcValue := c.Param("calcValue")
        return c.String(http.StatusOK, "result:" + calcValue)
    }
}
