package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type Address struct {
	ID       int    `json:"id"`
	PostCode string `json:"postcode"`
	HouseNo  string `json:"house_no"`
	RoadNo   string `json:"road_no"`
	RoadName string `json:"road_name"`
	Area     string `json:"area"`
	District string `json:"district"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
}

func getAddresses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := dbConn()

	defer db.Close()

	rows, err := db.Query("SELECT id, postcode, house_no, road_no, road_name, area, district FROM addresses")

	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()

	var addresses = []Address{}
	for rows.Next() {
		var (
			id       int
			postCode string
			houseNo  string
			roadNo   string
			roadName string
			area     string
			district string
		)
		err = rows.Scan(&id, &postCode, &houseNo, &roadNo, &roadName, &area, &district)
		if err != nil {
			panic(err.Error())
		}

		addresses = append(addresses, Address{id, postCode, houseNo, roadNo, roadName, area, district})
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(addresses)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := dbConn()

	defer db.Close()

	rows, err := db.Query("SELECT id, first_name, email FROM users")

	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()

	var users = []User{}
	for rows.Next() {
		var (
			id        int
			firstName string
			email     string
		)
		err = rows.Scan(&id, &firstName, &email)
		if err != nil {
			panic(err.Error())
		}

		users = append(users, User{id, firstName, email})
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(users)
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "butterfly"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func handleRequests() {
	http.HandleFunc("/", getAddresses)
	http.HandleFunc("/users", getUsers)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
