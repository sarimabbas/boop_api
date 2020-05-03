package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
)

// DB global GORM var
var DB *gorm.DB

func initDB() {
	var err error
	// DEV MODE
	if err = godotenv.Load(".env"); err != nil {
		DB, err = gorm.Open("postgres",
			fmt.Sprintf("host=%s port=5432 user=%s dbname=%s sslmode=disable",
				os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_NAME")))
		if err != nil {
			panic(err)
		}
		return
	}
	// PROD MODE
	DB, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
}

func getDB() *gorm.DB {
	return DB
}

// User GORM model
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"type:varchar(100);unique_index" json:"email" validate:"required"`
	Username string `gorm:"type:varchar(100);unique_index" json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Boop GORM model
type Boop struct {
	gorm.Model
	FromUserID int
	ToUserID   int
	FromUser   User `gorm:"foreignkey:FromUserID"`
	ToUser     User `gorm:"foreignkey:ToUserID"`
	Message    string
}

// CustomValidator GORM middleware
type CustomValidator struct {
	validator *validator.Validate
}

// Validate function for GORM
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
