package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Quote struct {
	Id     int
	Quote  string
	Author string
}

func main() {
	//load database credentials
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := 5432
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//Connect to the postgres database
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//create mux router
	router := mux.NewRouter()

	// /quotes route to get all quotes or create a new quote
	router.Handle("/quotes", GetQuotesHandler(db)).Methods("GET")
	router.Handle("/quotes", PostQuoteHandler(db)).Methods("POST")

	// run on port 8080
	http.ListenAndServe(":8080", router)
}
