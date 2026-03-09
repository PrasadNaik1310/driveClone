package db

import (
	"os"

	"log"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() error {
	dsn := os.Getenv("db_url")
	if dsn == "" {
		log.Fatalf("DB URL not found")
		return nil
	}

	/*	if strings.HasPrefix(dsn, "postgres://") {
		if(!strings.Contains(dsn,"connect_timeout")){
			u, err := url.Parse(dsn)
			if err == nil{
				q:= u.Query()
				q.Set("connect_timeout","10")
				u.RawQuery = q.Encode()
				dsn = u.String()

			}
		}
	}*/

	return nil
}
