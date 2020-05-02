package main

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// e.GET("/users/:username", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:username`
	id := c.Param("username")
	return c.String(http.StatusOK, id)
}

// e.POST("/users", createUser)
func createUser(c echo.Context) error {
	// name := c.FormValue("name")
	// email := c.FormValue("email")
	// username := c.FormValue("username")
	// password := c.FormValue("password")

	// unpack the above info from the JSON body into the model
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	// validate info is correct
	if err := c.Validate(u); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(hash)

	// attempt save into db
	db := getDB()
	if err = db.Create(&u).Error; err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusCreated, u)
}

// e.GET("/login", login)
func login(c echo.Context) error {
	// get username and password from form
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	// retrieve from db
	db := getDB()
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	// successfully logged in
	return c.JSON(http.StatusOK, user)
}

// e.GET("/boops", getBoops)
func getBoops(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	return c.String(http.StatusOK, username+" "+password)
}
