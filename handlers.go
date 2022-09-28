package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetQuotesHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//get all quotes and return to the user
		rows, err := db.Query(`SELECT "id", "quote", "author" FROM "quotes"`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []Quote
		for rows.Next() {
			var nextQuote Quote
			rows.Scan(&nextQuote.Id, &nextQuote.Quote, &nextQuote.Author)
			data = append(data, nextQuote)
		}

		payload, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(payload))
	}

	return http.HandlerFunc(fn)
}

func PostQuoteHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var q Quote
		var newQuote Quote

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Insert a new quote and return it to the user
		row := db.QueryRow(`
			INSERT INTO quotes(quote, author)
			VALUES($1, $2)
			RETURNING id, quote, author
		`, q.Quote, q.Author)
		row.Scan(&newQuote.Id, &newQuote.Quote, &newQuote.Author)

		payload, _ := json.Marshal(newQuote)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(payload))
	}

	return http.HandlerFunc(fn)
}
