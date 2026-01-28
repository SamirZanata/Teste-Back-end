package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	DB         DBConfig
	FreteRapido FreteRapidoConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type FreteRapidoConfig struct {
	BaseURL        string
	Token          string
	PlatformCode   string
	ShipperCNPJ    string
	DispatcherCEP  string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "quote_api"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		FreteRapido: FreteRapidoConfig{
			BaseURL:       getEnv("FRETE_RAPIDO_BASE_URL", "https://sp.freterapido.com"),
			Token:         getEnv("FRETE_RAPIDO_TOKEN", "1d52a9b6b78cf07b08586152459a5c90"),
			PlatformCode:  getEnv("FRETE_RAPIDO_PLATFORM_CODE", "5AKVkHqCn"),
			ShipperCNPJ:   getEnv("FRETE_RAPIDO_SHIPPER_CNPJ", "25438296000158"),
			DispatcherCEP: getEnv("FRETE_RAPIDO_DISPATCHER_CEP", "29161376"),
		},
	}
}

func (c *DBConfig) DSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.DBName + "?sslmode=" + c.SSLMode
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func GetIntEnv(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}
