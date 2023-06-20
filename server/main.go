package main

import (
	db "app/db"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"fmt"
)

type Ride struct {
	Id       string `json:"id"`
	CarId    string `json:"car_id"`
	Location string `json:"location"`
	Path     string `json:"path"`
}

type Customer struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Active   bool   `json:"active"`
	Location string `json:"location"`
}

func getRides(w http.ResponseWriter, req *http.Request) {
	//set cors for specific origin
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	rows, err := db.Connection.Query("SELECT * FROM rides")
	if err != nil {
		http.Error(w, "Failed to get rides: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rides []Ride

	for rows.Next() {
		var ride Ride
		rows.Scan(&ride.Id, &ride.CarId, &ride.Location, &ride.Path)
		rides = append(rides, ride)
	}

	ridesBytes, _ := json.MarshalIndent(rides, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(ridesBytes)
}

func getCustomers(w http.ResponseWriter, req *http.Request) {
	//set cors for specific origin
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	rows, err := db.Connection.Query("SELECT * FROM customers where active = true")
	if err != nil {
		http.Error(w, "Failed to get customers: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var customer Customer
		rows.Scan(&customer.Id, &customer.Name, &customer.Active, &customer.Location)
		customers = append(customers, customer)
	}

	ridesBytes, _ := json.MarshalIndent(customers, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(ridesBytes)
}


func main() {
	db.InitDB()
	defer db.Connection.Close()

	http.Handle("/", http.FileServer(http.Dir("../frontend/build")))
	http.HandleFunc("/rides", getRides)
	http.HandleFunc("/customers", getCustomers)

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	serverEnv := os.Getenv("SERVER_ENV")

	if serverEnv == "DEV" {
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else if serverEnv == "PROD" {
		// not implemented
		fmt.Println("PROD server not implemented yet.")
	}
}