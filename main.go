package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	var err error
	// Construct DSN from environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Open connection (this validates arguments, does not connect yet)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	fmt.Println("Successfully connected to MySQL")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	// Execute a simple query
	var dbMessage string
	err := db.QueryRow("SELECT 'Hello from MySQL!'").Scan(&dbMessage)
	if err != nil {
		http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Web Server says: Hello, World!\nDatabase says: %s", dbMessage)
}

func main() {
	// Initialize Database Connection
	initDB()
	// Ensure DB is closed when main exits
	defer db.Close()

	http.HandleFunc("/", helloHandler)

	fmt.Println("Server starting on :8085...")
	if err := http.ListenAndServe(":8085", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}