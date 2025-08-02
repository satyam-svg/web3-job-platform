package config

import (
	"log"
	"os"
)

type Config struct {
	GeminiAPIKey        string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
}

var AppConfig *Config

func LoadConfig() {
	geminiKey := os.Getenv("GEMINI_API_KEY")
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	cloudKey := os.Getenv("CLOUDINARY_API_KEY")
	cloudSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Fail fast if any required secret is missing
	if geminiKey == "" {
		log.Fatal("❌ GEMINI_API_KEY environment variable not set")
	}
	if cloudName == "" || cloudKey == "" || cloudSecret == "" {
		log.Fatal("❌ Cloudinary environment variables not set")
	}

	AppConfig = &Config{
		GeminiAPIKey:        geminiKey,
		CloudinaryCloudName: cloudName,
		CloudinaryAPIKey:    cloudKey,
		CloudinaryAPISecret: cloudSecret,
	}
}
