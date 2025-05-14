package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DSN string
}

func Load() *Config {
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("No config.env file found, using env vars")
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	//dns := fmt.Sprintf("host=%s port =%s user=%s password%s dbname=%s sslmode=disable",
	dns := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		db,
	)

	return &Config{DSN: dns}
}
