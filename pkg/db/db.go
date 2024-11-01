package db

import (
	"fmt"
	"log"

	"github.com/shivaraj-shanthaiah/user-management/config"
	"github.com/shivaraj-shanthaiah/user-management/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(config *config.Config) *gorm.DB {
	host := config.DBHost
	user := config.DBUser
	password := config.DBPassword
	dbname := config.DBDatabase
	port := config.DBPort
	sslmode := config.DBSslmode

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("connection to postgres database failed: ", err)
	}
	log.Printf("Successfully connected to postgres Host:%s,  DBName:%s, Port:%s", host, dbname, port)

	err = DB.AutoMigrate(
		&model.User{},
	)

	if err != nil {
		log.Printf("error while migrating %v \n", err.Error())
		return nil
	}
	return DB
}
