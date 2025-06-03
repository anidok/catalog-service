package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secret = "some-key" // Same secret as in Kong config

func generateToken() {
	claims := jwt.MapClaims{
		"iss": "kong-jwt-auth",
		"exp": time.Now().Add(time.Hour).Unix(),
		"sub": "user125",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Failed to sign token:", err)
		os.Exit(1)
	}
	fmt.Println("Generated JWT Token:")
	fmt.Println(signed)
}

func verifyToken(tokenStr string) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Token verification failed:", err)
		os.Exit(1)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token is valid. Decoded payload:")
		for k, v := range claims {
			fmt.Printf("  %s: %v\n", k, v)
		}
	} else {
		fmt.Println("Token is invalid.")
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go generate")
		fmt.Println("  go run main.go verify <token>")
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
	case "generate":
		generateToken()
	case "verify":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go verify <token>")
			os.Exit(1)
		}
		verifyToken(os.Args[2])
	default:
		fmt.Println("Usage:")
		fmt.Println("  go run main.go generate")
		fmt.Println("  go run main.go verify <token>")
		os.Exit(1)
	}
}
