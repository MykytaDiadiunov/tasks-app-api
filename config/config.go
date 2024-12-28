package config

import (
	"os"
)

type Configuration struct {
	DatabaseName        string
	DatabaseHost        string
	DatabasePort        string
	DatabaseUser        string
	DatabasePassword    string
	JwtSecret           string
	MigrateToVersion    string
	MigrationLocation   string
	LoggerLevel         string
	SmtpHost            string
	SmtpPort            string
	WorkGmail           string
	WorkGmailPassword   string
	CloudinaryNameKey   string
	CloudinaryApiKey    string
	CloudinarySecretKey string
}

func GetConfiguration() Configuration {
	return Configuration{
		DatabaseName:        getOrDefault("DB_NAME", "restapi_dev"),
		DatabaseHost:        getOrDefault("DB_HOST", "127.0.0.1"),
		DatabasePort:        getOrDefault("DB_PORT", "5432"),
		DatabaseUser:        getOrDefault("DB_USER", "postgres"),
		DatabasePassword:    getOrDefault("DB_PASSWORD", "postgres"),
		JwtSecret:           getOrDefault("JWT_SECRET", "1234567890"),
		MigrateToVersion:    getOrDefault("MIGRATE", "latest"),
		MigrationLocation:   getOrDefault("MIGRATION_LOCATION", "migrations"),
		LoggerLevel:         getOrDefault("LOGGER_LEVEL", "dev"),
		SmtpHost:            getOrDefault("SMTP_HOST", "smtp.gmail.com"),
		SmtpPort:            getOrDefault("SMTP_PORT", "587"),
		WorkGmail:           getOrDefault("WORK_GMAIL", "test@gmail.com"),
		WorkGmailPassword:   getOrDefault("WORK_GMAIL_PASSWORD", "test_password"),
		CloudinaryNameKey:   getOrDefault("CLOUDINARY_NAME_KEY", ""),
		CloudinaryApiKey:    getOrDefault("CLOUDINARY_API_KEY", ""),
		CloudinarySecretKey: getOrDefault("CLOUDINARY_SECRET_KEY", ""),
	}
}

func getOrDefault(key, defaultVal string) string {
	env, set := os.LookupEnv(key)
	if !set {
		return defaultVal
	}
	return env
}
