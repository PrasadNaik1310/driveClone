package db

import (
	"net/url"
	"os"
	"strings"
	"time"

	"log"

	"github.com/PrasadNaik1310/driveClone/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDb() error {
	dsn := os.Getenv("db_url")
	if dsn == "" {
		log.Fatalf("DB URL not found coming from dbinit")
		return nil
	}

	if strings.HasPrefix(dsn, "postgres://") {
		if !strings.Contains(dsn, "connect_timeout") {
			u, err := url.Parse(dsn)
			if err == nil {
				q := u.Query()
				q.Set("connect_timeout", "10")
				u.RawQuery = q.Encode()
				dsn = u.String()

			}
		}
	} else if !strings.Contains(dsn, "connect_timeout") {
		// DSN format (key=value pairs)
		dsn += "&X connect_timeout=10"
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("Nahi ho rha connect ")
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}
	DB = database
	log.Print("Connected to database successfully")

	//configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetMaxIdleConns(102)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	//auto-migrate
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Print("Failed to AtuoMigrate Users")
		return err
	}
	log.Printf("user migrated")
	if err := DB.AutoMigrate(&models.Folder{}); err != nil {
		log.Print("failed to Automigrate Folder table")
		return err
	}
	log.Printf("Folder migrated")
	if err := DB.AutoMigrate(&models.File{}); err != nil {
		log.Print("failed to Automigrate File table")
		return err
	}
	log.Printf("File table migrated")
	log.Printf("DB Chaluuuuuuuuuu")

	return nil
}
