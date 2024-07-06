package base

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load() // 👈 load .env file
	if err != nil {
		log.Fatal(err)
	}
}
