package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func checkEnv(key string) {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s empty in env", key)
	}
}

func LoadEnvs(strict bool) {
	err := godotenv.Load()
	if err == nil {
		log.Print("Using .env file")
	}
	if strict {
		checkEnv("API_URL")
		checkEnv("API_KEY")
		checkEnv("JWT_MASTER_SECRET")
		checkEnv("DB_URL")
	}
}
