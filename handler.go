package main

import (
	"fmt"

	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	
)

var db *gorm.DB
var err error

type Model struct {
	ID        uint `gorm:"primary_key"` 
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	DeletedAt *time.Time `json:"deleted"`
}

type Product struct {
	ProductID int `gorm:"primary_key"`
	ProductTitle       string
	ProductDescription string
	ProductImages       []ProductImage `gorm:"ForeignKey:Proid"`
	Brand              string
	Category           string
	Price              string
	ProductStatus      string
	Subcategory        string
}



type ProductImage struct {
	ImageID int `gorm:"primary_key"`
	Image string
	Proid int
}

type Orders struct {
	gorm.Model
	ProductID int
	UserID  uint
}

type User struct {
	gorm.Model
	Email        string  `gorm:"type:varchar(100);unique"`
	Password     string  `json:"-"`
	ProfileImage string `json:"avatar"`
	Role         string `json:"role"`
}

type JwtToken struct {
	Token string
}



func InitialMigration() {
	db, err = gorm.Open("sqlite3", "ecommerce.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()
	db.AutoMigrate(&Product{}, &User{}, &Orders{},&ProductImage{})
}




