package database

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func GetInstance() (*gorm.DB, error) {
	var err error

	once.Do(func() {
		// Load env once
		err = godotenv.Load()
		if err != nil {
			log.Println("Warning: error loading .env file, continuing with environment variables")
		}

		driver := os.Getenv("DB_Driver")
		if driver == "" {
			err = fmt.Errorf("DB_Driver environment variable not set")
			return
		}

		switch driver {
		case "mysql":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				os.Getenv("DB_User"),
				os.Getenv("DB_Password"),
				os.Getenv("DB_HOST"),
				os.Getenv("DB_Port"),
				os.Getenv("DB_Name"),
			)
			dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		case "postgresql", "postgres":
			sslmode := os.Getenv("DB_SSLMODE")
			if sslmode == "" {
				sslmode = "disable"
			}
			dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
				os.Getenv("DB_HOST"),
				os.Getenv("DB_Port"),
				os.Getenv("DB_User"),
				os.Getenv("DB_Password"),
				os.Getenv("DB_Name"),
				sslmode,
				os.Getenv("DB_TimeZone"),
			)
			dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		default:
			err = fmt.Errorf("unsupported DB_Driver: %s", driver)
		}

		if err != nil {
			log.Printf("Error connecting to database: %v", err)
		} else {
			log.Println("Database connected successfully")
		}
	})

	return dbInstance, err
}
