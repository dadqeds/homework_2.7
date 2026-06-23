package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./users.db"

func authenticate(db *sql.DB, email, password string) string {
	var stored string
	err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&stored)
	if err == sql.ErrNoRows {
		return "пользователь не найден"
	}
	if err != nil {
		log.Fatal("query:", err)
	}
	if stored != hashPassword(password) {
		return "пароль неверный"
	}
	return "пароль верный"
}

func main() {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("open db:", err)
	}
	defer db.Close()

	if len(os.Args) == 1 {
		seedDB(db)
		fmt.Println("\nБД заполнена. Запустите с аргументами: go run . <email> <password>")
		return
	}

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Использование: %s <email> <password>\n", os.Args[0])
		os.Exit(1)
	}

	email := os.Args[1]
	password := os.Args[2]

	result := authenticate(db, email, password)
	fmt.Println(result)
}
