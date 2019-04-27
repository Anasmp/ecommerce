package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter()
	// userhandlers
	myRouter.HandleFunc("/registerUser", RegisterUser).Methods("POST")
	myRouter.HandleFunc("/loginUser", loginUser).Methods("POST")
	myRouter.HandleFunc("/listallUsers", AllCustomers).Methods("POST")
	//producthandlers
	myRouter.HandleFunc("/addproduct", AddProduct).Methods("POST")
	// myRouter.HandleFunc("/updateproduct", UpdateProduct).Methods("POST")
	myRouter.HandleFunc("/deleteproduct", DeleteProduct).Methods("POST")
	myRouter.HandleFunc("/products", AllProducts).Methods("POST")
	myRouter.HandleFunc("/searchProduct", SearchProduct).Methods("POST")
	//order handlers
	myRouter.HandleFunc("/submitorder", SubmitOrder).Methods("POST")
	myRouter.HandleFunc("/getorders", getOrders).Methods("POST")

	myRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func main() {
	fmt.Println("Server started at http://localhost:8000")
	InitialMigration()
	handleRequests()
}
