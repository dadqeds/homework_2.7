package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	_ "github.com/mattn/go-sqlite3"
)

var firstNames = []string{"alice", "bob", "carol", "dave", "eve", "frank", "grace", "heidi"}
var domains = []string{"gmail.com", "yahoo.com", "mail.ru", "example.com"}
var rawPasswords = []string{"qwerty123", "letmein", "hunter2", "pass1234", "secret!", "abc@xyz", "P@ssw0rd", "sunshine"}

func hashPassword(pw string) string {
	h := sha256.Sum256([]byte(pw))
	return fmt.Sprintf("%x", h)
}

func seedDB(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		email    TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("create table:", err)
	}

	stmt, err := db.Prepare("INSERT OR IGNORE INTO users (email, password) VALUES (?, ?)")
	if err != nil {
		log.Fatal("prepare:", err)
	}
	defer stmt.Close()

	fmt.Println("Seeding users:")
	for i := 0; i < 10; i++ {
		name := firstNames[rand.Intn(len(firstNames))] + fmt.Sprintf("%d", rand.Intn(100))
		email := name + "@" + domains[rand.Intn(len(domains))]
		pw := rawPasswords[rand.Intn(len(rawPasswords))]
		_, err := stmt.Exec(email, hashPassword(pw))
		if err != nil {
			log.Printf("  skip %s: %v", email, err)
		} else {
			fmt.Printf("  %-35s  password: %s\n", email, pw)
		}
	}
}
