package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitializeDatabase() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	DB, err = gorm.Open(postgres.Open(psqlInfo), gormConfig())

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database Successfully connected!")
	}

	db, errDb := DB.DB()
	if errDb != nil {
		panic(errDb)
	}

	setMaxIdle, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONS"))
	setMaxOpen, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONS"))
	db.SetMaxIdleConns(setMaxIdle)
	db.SetMaxOpenConns(setMaxOpen)
	migrateDatabase(DB)
}

func gormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
			},
		),
	}
}

func migrateDatabase(db *gorm.DB) {
	if err := db.AutoMigrate(&Name{}); err != nil {
		panic(err)
	}
}
