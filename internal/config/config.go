package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	//Dev env
	Environment string

	//Server configurations
	PublicHost string
	Port       string

	//Database configurations
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string

	//JWT Secret
	JWTSecret string
}

var Envs = initConfig()

func initConfig() Config {
	//load the .env file
	godotenv.Load()
	// Load environment variables from .env file if it exists
	return Config{
		PublicHost:  getEnv("PUBLIC_HOST", "http://localhost:8082"),
		Port:        getEnv("PORT", "8082"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", "phantom0627"),
		DBAddress:   fmt.Sprintf("%s:%s", getEnv("DB_ADDRESS", "localhost"), getEnv("DB_PORT", "3306")),
		DBName:      getEnv("DB_NAME", "authentication_service"),
		JWTSecret:   getEnv("JWT_SECRET", "2XLkQROFwhxU4myufKGdYISNJB87bgzP9vCiEtHpWV0cMqATjan3solZr165De"),
		Environment: getEnv("ENV", "dev"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
