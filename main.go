package main

import (
	// "encoding/json"
	"github.com/go-playground/validator/v10"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	initDB()
	db := getDB()
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})

	// Init router
	e := echo.New()

	// add validation
	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/users/:username", getUser)
	e.POST("/users", createUser)
	e.GET("/boops", getBoops)
	e.GET("/login", login)

	// e.POST("/user", func(c echo.Context) error {
	// 	return
	// })

	// 	e.POST("/users", saveUser)
	// e.GET("/users/:id", getUser)
	// e.PUT("/users/:id", updateUser)
	// e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}
