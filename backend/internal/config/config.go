package config

import (
	"os"
	"strings"
)

type Settings struct {
	DatabaseURL    string
	JWTSecret      string
	JWTExpireHours int
	CORSOrigins    []string
}

func Load() Settings {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "sqlite:///./taskflow.db"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}

	originsRaw := os.Getenv("CORS_ORIGINS")
	if originsRaw == "" {
		originsRaw = "http://localhost:8000,http://localhost:5500"
	}
	var origins []string
	for _, o := range strings.Split(originsRaw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}

	return Settings{
		DatabaseURL:    dbURL,
		JWTSecret:      secret,
		JWTExpireHours: 24,
		CORSOrigins:    origins,
	}
}
