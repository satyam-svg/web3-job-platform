package config

import (
	"log"
	"os"

	"github.com/satyam-svg/resume-parser/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = "resume.db" // default name
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to SQLite: %v", err)
	}

	log.Println("✅ Connected to SQLite DB")

	// Auto-migrate tables
	log.Println("🔄 Starting database migration...")

	if err := DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("❌ User table migration failed: %v", err)
	}
	log.Println("✅ User table migrated successfully")

	if err := DB.AutoMigrate(&model.Education{}); err != nil {
		log.Fatalf("❌ Education table migration failed: %v", err)
	}
	log.Println("✅ Education table migrated successfully")

	if err := DB.AutoMigrate(&model.Experience{}); err != nil {
		log.Fatalf("❌ Experience table migration failed: %v", err)
	}
	log.Println("✅ Experience table migrated successfully")

	log.Println("✅ All tables migrated successfully")

	DB.AutoMigrate(&model.Job{})

	// Debug: List all tables
	var tables []string
	DB.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables)
	log.Println("📋 Database tables:", tables)

	// Optionally test tables
	testTables()

	return DB // ✅ Return DB here
}

func testTables() {
	log.Println("🧪 Testing table creation...")

	// Test Education table
	testEducation := model.Education{
		Institution: "Test University",
		Location:    "Test City",
		Degree:      "Test Degree",
		GPA:         "3.5",
		Years:       "2020-2024",
	}

	if err := DB.Create(&testEducation).Error; err != nil {
		log.Printf("❌ Education table test failed: %v", err)
	} else {
		log.Println("✅ Education table test passed")
		DB.Delete(&testEducation)
	}

	// Test Experience table
	testExperience := model.Experience{
		Company:     "Test Company",
		Location:    "Test City",
		Title:       "Test Position",
		Years:       "2022-2024",
		Description: "Test description",
	}

	if err := DB.Create(&testExperience).Error; err != nil {
		log.Printf("❌ Experience table test failed: %v", err)
	} else {
		log.Println("✅ Experience table test passed")
		DB.Delete(&testExperience)
	}
}
