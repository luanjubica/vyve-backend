#!/bin/bash
# Generate a bcrypt hash for a password
# Usage: ./generate_hash.sh <password>

PASSWORD="${1:-password123}"

cat > /tmp/gen_hash.go << 'EOF'
package main

import (
	"fmt"
	"log"
	"os"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := os.Args[1]
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(hashedPassword))
}
EOF

cd /tmp && go run gen_hash.go "$PASSWORD"
