package config

import (
	"log"

	"github.com/joho/godotenv"
)

func Load() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: godotenv.Load() failed: %v", err)
		log.Println("Will attempt to use system environment variables instead")
	} else {
		log.Println("Successfully loaded .env file")
	}
	return nil
}
