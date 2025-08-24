package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	MongoURI string
	DbName   string
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func LoadConfig() Config {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found, using system env variables")
	}

	return Config{
		Port:     getEnv("PORT", "8080"),
		MongoURI: getEnv("MONGO_URI", ""),
		DbName:   getEnv("DB_NAME", "jobster"),
	}
}
