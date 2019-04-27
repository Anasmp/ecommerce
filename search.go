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

func SearchProduct(w http.ResponseWriter, r *http.Request) {
	
	tokenmain := r.FormValue("token")
	searchkey := r.FormValue("key")

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

		var product []Product
		// db.Where("product_title LIKE ? AND  product_status = ?", "%"+"Mac"+"%", "publish").Limit(10).Order("id desc").Find(&product)
		db.Where("product_title LIKE ? AND  product_status = ?", "%"+searchkey+"%", "publish").Preload("ProductImages").Limit(10).Order("product_id desc").Find(&product)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)

	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"invalid token"}`)
	}
}