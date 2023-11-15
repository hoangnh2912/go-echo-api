package db

import (
	"fmt"
	"os"

	"go-echo-api/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetPSQLInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+" password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), 5432, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
}

func New() *gorm.DB {
	fmt.Println(GetPSQLInfo())
	postgres.Open(GetPSQLInfo())
	db, err := gorm.Open(postgres.Open(GetPSQLInfo()))
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	// db.DB().SetMaxIdleConns(3)
	dbConfig, err := db.DB()
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	dbConfig.SetMaxIdleConns(3)
	db.Logger.LogMode(logger.Info)
	return db
}

func TestDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(GetPSQLInfo()))
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	dbConfig, err := db.DB()
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	dbConfig.SetMaxIdleConns(3)
	db.Logger.LogMode(logger.Info)
	return db
}

func DropTestDB() error {
	db, err := gorm.Open(postgres.Open(GetPSQLInfo()))
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	db.Exec("DROP DATABASE IF EXISTS postgres")
	return nil

}

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
