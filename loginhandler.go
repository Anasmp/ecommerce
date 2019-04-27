package main

import(
	"encoding/json"
	"io/ioutil"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	r.ParseMultipartForm(10 << 20)
	file, handler, errprofile := r.FormFile("profileImage")


	if len(email) > 5 && len(password) > 6 {

		db, err = gorm.Open("sqlite3", "ecommerce.db")
		if err != nil {
			panic("Could not connect to the database")
		}
		defer db.Close()

		var user []User

		db.Where("email = ?", email).First(&user)

		if len(user) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"message":"user already exist"}`)
		} else {
	

			hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
			if err != nil {
				panic("password hashing failed")
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"email":    email,
				"password": string(hashpassword),
			})
			tokenString, error := token.SignedString([]byte("im-codedady-supertoken"))
			if error != nil {
				fmt.Println(error)
			}

			if(errprofile == nil){
					
					tempFile, err := ioutil.TempFile("assets/profile", "user-*.png")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Printf("Uploaded File: %+v\n", handler.Filename)
					defer tempFile.Close()
					fileBytes, err := ioutil.ReadAll(file)
					if err != nil {
						fmt.Println(err)
					}
		
					tempFile.Write(fileBytes)
					filepath := tempFile.Name()

					db.Create(&User{Email: email, Password: string(hashpassword), ProfileImage: filepath, Role: "customer"})
			}else{
				    db.Create(&User{Email: email, Password: string(hashpassword), Role: "customer"})
			}
			
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(JwtToken{Token: tokenString})

		}}else {
		// fmt.Fprintf(w, "check your password or username")
	    w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"your password or username too short"}`)
	}

}

func loginUser(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(email) > 5 && len(password) > 6 {

		db, err = gorm.Open("sqlite3", "ecommerce.db")
		if err != nil {
			panic("Could not connect to the database")
		}
		defer db.Close()

		var user []User

		db.Where("email = ? ", email).First(&user)

		if len(user) > 0 {
			match := CheckPasswordHash(password, user[0].Password)
			if match == true {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"email":    email,
					"password": user[0].Password,
				})
				tokenString, error := token.SignedString([]byte("im-codedady-supertoken"))
				if error != nil {
					fmt.Println(error)
				}
				json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"message":"enter a valid password"}`)
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"message":"no user found"}`)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"check your password or username"}`)
	}

}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}