package main

import "golang.org/x/crypto/bcrypt"

func main() {

	hash := "$2a$10$rK2Ia/SRslB0c7GTTsHKqewhFoccS/Q/189UHCkiT6HIMKg.lHgHi"

	password := "bellow"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		panic(err)
	}

	hash = "$2a$10$g6efNvur33Ya9lV1Fo3ClekXylbUgv2JEl.9.21vCFRPEoBopk6k."

	password = "lucky"

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		panic(err)
	}
}
