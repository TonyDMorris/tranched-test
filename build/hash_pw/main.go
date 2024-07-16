package main

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	password := os.Args[1]

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	log.Println(string(hashedPassword))
}
