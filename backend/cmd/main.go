package main

import (
	"github.com/williamf6894/VB-Events/internal/db"
)

func main() {
	_, err := db.InitDB()
	if err != nil {
		panic("failed to connect to database")
	}
}
