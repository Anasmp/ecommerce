package main

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"io"
	"fmt"
	"strconv"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mitchellh/mapstructure"
)


//add product
func AddProduct(w http.ResponseWriter, r *http.Request) {

	tokenmain := r.FormValue("token")
	productDescription := r.FormValue("productDescription")
	productTitle := r.FormValue("productTitle")
	brandName := r.FormValue("brand")
	categoryName := r.FormValue("category")
	productPrice := r.FormValue("price")
	productStatus := r.FormValue("status")
	subcategory := r.FormValue("subcategory")

	if len(tokenmain) > 1 && len(productTitle) > 1 {

		var arr []ProductImage

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
			// json.NewEncoder(w).Encode(user)
			// fmt.Println(user.Email)

			db, err = gorm.Open("sqlite3", "ecommerce.db")
			if err != nil {
				panic("Could not connect to the database")
			}
			defer db.Close()

		

			err := r.ParseMultipartForm(100000)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		
		
			m := r.MultipartForm
		
			
			files := m.File["images"]
			for i, _ := range files {
				
				file, err := files[i].Open()
				defer file.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				
				dst, err := ioutil.TempFile("assets/product", "products-*.png")
				defer dst.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				
				if _, err := io.Copy(dst, file); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
		
			
				a := ProductImage{Image: dst.Name()}
				arr = append(arr,a)
				
			}
			

			poduct1 := Product{ProductTitle: productTitle, ProductDescription: productDescription,ProductImages:arr, Brand: brandName, Category: categoryName, Price: productPrice, ProductStatus: productStatus, Subcategory: subcategory}
			db.Create(&poduct1)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message":"product added successfully"}`)

		} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"invalid token"}`)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"invalid token"}`)
	}
}


// updateproduct

func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	tokenmain := r.FormValue("token")
	productid := r.FormValue("productid")

	if len(tokenmain) > 1 && len(productid) >= 1 {

		token, _ := jwt.Parse(tokenmain,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte("im-codedady-supertoken"), nil
			})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var product Product
			var user User
			mapstructure.Decode(claims, &user)

			db, err = gorm.Open("sqlite3", "ecommerce.db")
			if err != nil {
				panic("Could not connect to the database")
			}
			defer db.Close()
			intproductid, _ := strconv.Atoi(productid)
			if db.Where(&Product{ProductID: intproductid}).Find(&product).RecordNotFound() {
				// record not found
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"product not found in database"}`)
			  }else{
				db.Delete(&product)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"product deleted successfully"}`)
			  }
			
			

		

		} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"invalid token"}`)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"invalid token"}`)
	}
}

















func AllProducts(w http.ResponseWriter, r *http.Request) {

	tokenmain := r.FormValue("token")
	 offset := r.FormValue("page")
	status := r.FormValue("status")
	category := r.FormValue("category")
	subcategory := r.FormValue("subcategory")

	if len(tokenmain) > 4 {

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

			offsetValue, err := strconv.ParseInt(offset, 10, 64)
			if err == nil {
				fmt.Println(offsetValue)
			}
			// offsetAdd := offsetValue * 10
			// offsetAddnew := strconv.Itoa(int(offsetAdd))

			var product []Product
			// product := &Product{}
			// db.Debug().Where("product_title=?","mac book").Preload("ProductImages").Find(&product) //db.Debug().Where("customer_name=?","John").Preload("Contacts").Find(&customers)
			// json.NewEncoder(w).Encode(product)

			if len(category) > 1 && len(status) > 1 && len(offset) >= 1 && len(subcategory) >= 1 {
				db.Where(&Product{ProductStatus: status, Category: category, Subcategory: subcategory}).Preload("ProductImages").Offset(offset).Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) > 1 && len(status) > 1 && len(offset) >= 1 && len(subcategory) < 1 {
				db.Where(&Product{ProductStatus: status, Category: category}).Preload("ProductImages").Offset(offset).Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) > 1 && len(status) > 1 && len(offset) <= 1 && len(subcategory) > 1 {
				db.Where(&Product{ProductStatus: status, Category: category, Subcategory: subcategory}).Preload("ProductImages").Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) > 1 && len(status) > 1 && len(offset) <= 1 && len(subcategory) < 1 {
				db.Where(&Product{ProductStatus: status, Category: category}).Preload("ProductImages").Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) > 1 && len(status) < 1 && len(offset) >= 1 && len(subcategory) > 1 {
				db.Where(&Product{Category: category, Subcategory: subcategory}).Preload("ProductImages").Offset(offset).Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) > 1 && len(status) < 1 && len(offset) >= 1 && len(subcategory) < 1 {
				db.Where(&Product{Category: category}).Preload("ProductImages").Offset(offset).Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) > 1 && len(status) < 1 && len(offset) <= 1 && len(subcategory) > 1 {
				db.Where(&Product{Category: category, Subcategory: subcategory}).Preload("ProductImages").Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) > 1 && len(status) < 1 && len(offset) <= 1 && len(subcategory) < 1 {
				db.Where(&Product{ProductStatus: status, Category: category, Subcategory: subcategory}).Preload("ProductImages").Offset(offset).Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) < 1 && len(status) > 1 && len(offset) >= 1 && len(subcategory) > 1 {
				db.Where(&Product{ProductStatus: status, Subcategory: subcategory}).Preload("ProductImages").Offset(offset).Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) < 1 && len(status) > 1 && len(offset) >= 1 && len(subcategory) < 1 {
				db.Where(&Product{ProductStatus: status}).Preload("ProductImages").Offset(offset).Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) < 1 && len(status) > 1 && len(offset) <= 1 && len(subcategory) > 1 {
				db.Where(&Product{ProductStatus: status, Subcategory: subcategory}).Preload("ProductImages").Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) < 1 && len(status) > 1 && len(offset) <= 1 && len(subcategory) < 1 {

				db.Debug().Where(&Product{ProductStatus: status}).Limit(10).Order("product_id desc").Preload("ProductImages").Find(&product)

			} else if len(category) < 1 && len(status) < 1 && len(offset) >= 1 && len(subcategory) > 1 {
				db.Where(&Product{Subcategory: subcategory}).Preload("ProductImages").Offset(offset).Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) < 1 && len(status) < 1 && len(offset) >= 1 && len(subcategory) < 1 {
				db.Preload("ProductImages").Offset(offset).Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) < 1 && len(status) < 1 && len(offset) <= 1 && len(subcategory) > 1 {
				db.Where(&Product{Subcategory: subcategory}).Preload("ProductImages").Limit(10).Order("product_id desc").Find(&product)
			} else if len(category) < 1 && len(status) < 1 && len(offset) <= 1 && len(subcategory) < 1 {
				db.Preload("ProductImages").Limit(10).Order("product_id desc").Find(&product)
			}

			w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(product)

		} else {
			w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"invalid token"}`)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"invalid request"}`)
	}
}


