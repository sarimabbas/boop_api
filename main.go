package main

import (
	// "encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// setup db
	initDB()
	db := getDB()
	defer db.Close()

	// migrate schema
	db.AutoMigrate(&User{}, &Boop{})

	// Init router
	e := echo.New()

	// add validation
	e.Validator = &CustomValidator{validator: validator.New()}

	// add JWT middleware
	r := e.Group("")
	r.Use(middleware.JWT([]byte("secret")))

	// add routes
	e.GET("/", welcome)
	e.GET("/users/:username", getUser)
	e.POST("/users", createUser)
	e.GET("/boops", getBoops)
	e.POST("/login", login)
	r.GET("/test", testToken)
	r.POST("/boops", sendBoop)

	// start server
	e.Logger.Fatal(e.Start(":1323"))
}
