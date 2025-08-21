package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run check_password.go <password> <bcrypt-hash>")
	}

	password := os.Args[1]
	hash := os.Args[2]

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Password comparison failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Password matches hash!")
}
