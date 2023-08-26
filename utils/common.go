package utils

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetEnvKey : return Environment value
func GetEnvKey(env string) string {
	key, ok := os.LookupEnv(env)
	if !ok {
		log.Println(env + " environment variable required but not set")
		os.Exit(1)
	}
	return key
}

//CheckEnv : Check if the environment variables are set
func CheckEnv() error {
	godotenv.Load()
	_, ok := os.LookupEnv("DB_URL")
	if !ok {
		return errors.New("DB_URL environment variable required but not set")
	}
	_, ok = os.LookupEnv("JWTKEY")
	if !ok {
		return errors.New("JWTKEY environment variable required but not set")
	}
	return nil
}
