package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Person struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "database"
	dbname   = "postgres"
)

func OpenConnection() *sql.DB {
	fmt.Println("Trying opening the databasee")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("error while opening the database")
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("error while pinging the databasee")
		panic(err)
	}
	fmt.Println("Database opened")
	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM person")
	if err != nil {
		log.Fatal(err)
	}

	var people []Person

	for rows.Next() {
		var person Person
		rows.Scan(&person.Name, &person.Nickname)
		people = append(people, person)
	}

	peopleBytes, _ := json.MarshalIndent(people, "T", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POSTHandler hit")
	db := OpenConnection()

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Println("TEST Error decoding the received http request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO person (name, nickname) VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, p.Name, p.Nickname)
	if err != nil {
		fmt.Println("TEST EXEC command failed", p.Name, p.Nickname)
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func main() {
	http.HandleFunc("/", GETHandler)
	http.HandleFunc("/insert", POSTHandler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
