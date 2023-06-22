package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "c:/users/lizyu/Moonshine-backend/db/users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建用户表
	createTable()

	http.HandleFunc("/register", registerHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTable() {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			password TEXT
		);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	insertSQL := `
		INSERT INTO users (username, password) VALUES (?, ?);
	`

	_, err = db.Exec(insertSQL, user.Username, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
