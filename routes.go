package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

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

// e.POST("/login", login)
func login(c echo.Context) error {
	// get username and password from form
	form := new(User)
	if err := c.Bind(form); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	// retrieve from db
	db := getDB()
	var user User
	if err := db.Where("username = ?", form.Username).First(&user).Error; err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	// * YAY! All checks cleared

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

// Validate your JWT. If valid, it will be dumped back, else a return
func testToken(c echo.Context) error {
	token := c.Get("user")
	claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
	return c.JSON(http.StatusOK, "Welcome, "+claims["username"].(string))
}

// e.GET("/boops", getBoops)
func getBoops(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	return c.String(http.StatusOK, username+" "+password)
}

// e.POST("/boop", sendBoops)
func sendBoop(c echo.Context) error {
	// get receiver
	type Form struct {
		Username string `json:"username" validate:"required"`
		Message  string `json:"message"`
	}
	form := new(Form)
	if err := c.Bind(form); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	if err := c.Validate(form); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	db := getDB()
	var recipient User
	if err := db.Where("username = ?", form.Username).First(&recipient).Error; err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	// get sender
	token := c.Get("user")
	claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	var sender User
	if err := db.Where("username = ?", username).First(&sender).Error; err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	// create boop
	boop := Boop{FromUser: sender, ToUser: recipient, Message: form.Message}
	db.Create(&boop)
	return c.JSON(http.StatusCreated, boop)
}
