package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mitchellh/mapstructure"
)

func getOrders(w http.ResponseWriter, r *http.Request) {
	tokenmain := r.FormValue("token")

	db, err = gorm.Open("sqlite3", "ecommerce.db")
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	if len(tokenmain) > 1{
		token, _ := jwt.Parse(tokenmain,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte("im-codedady-supertoken"), nil
			})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var user User
			mapstructure.Decode(claims, &user)

			db, err = gorm.Open("sqlite3", "ecommerce.db")
			if err != nil {
				panic("Could not connect to the database")
			}
			defer db.Close()

			var user2 User
			
			var orders []Orders

			db.Where("email = ?",user.Email).First(&user2)
			db.Where("user_id = ?",user2.ID).Find(&orders)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(orders)
		
		}
	}
	
}