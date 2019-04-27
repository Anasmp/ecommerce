package main

import (
	"net/http"
	"fmt"
	"strconv"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mitchellh/mapstructure"
)

func SubmitOrder(w http.ResponseWriter, r *http.Request) {
	tokenmain := r.FormValue("token")
	productIdentity := r.FormValue("productid")

	db, err = gorm.Open("sqlite3", "ecommerce.db")
	if err != nil {
		panic("Could not connect to the database")
	}
	
	defer db.Close()

	if len(tokenmain) > 1 && len(productIdentity) > 0{
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
			productId,_:= strconv.Atoi(productIdentity)
			db.Where("email = ?",user.Email).First(&user2)
			db.Create(&Orders{UserID:user2.ID,ProductID:productId})

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message":"product added to cart"}`)
		
		}else{
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"message":"invalid token"}`)
		}
	}else{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"invalid request"}`)
	}
	
}