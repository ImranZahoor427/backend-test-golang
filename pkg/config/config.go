package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	MongoURI     string
	MongoDB      string
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
}

func Load() Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file; using default/fallback values.")
	}

	return Config{
		Port:         getEnv("PORT", "8080"),
		PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
		PostgresUser: getEnv("POSTGRES_USER", "postgres"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "password"),
		PostgresDB:   getEnv("POSTGRES_DB", "ledgerdb"),
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      getEnv("MONGO_DB", "testdb"),
		KafkaBrokers: splitAndTrim(getEnv("KAFKA_BROKERS", "localhost:9092")),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "transactions"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "transaction-consumer-group"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	var trimmed []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			trimmed = append(trimmed, p)
		}
	}
	return trimmed
}
