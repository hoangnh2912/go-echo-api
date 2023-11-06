package db

import (
	"fmt"
	"os"

	"go-echo-api/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const dialect = "postgres"

func GetPSQLInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+" password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), 5432, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
}

func New() *gorm.DB {
	db, err := gorm.Open(dialect, GetPSQLInfo())
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)
	return db
}

func TestDB() *gorm.DB {
	db, err := gorm.Open(dialect, GetPSQLInfo())
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(false)
	return db
}

func DropTestDB() error {
	db, err := gorm.Open(dialect, GetPSQLInfo())
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	db.Exec("DROP DATABASE IF EXISTS postgres")
	return nil

}

// TODO: err check
func AutoMigrate(db *gorm.DB) {
	fmt.Println("Starting Auto Migration")

	db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.Article{},
		&model.Comment{},
		&model.Tag{},
	)
}
