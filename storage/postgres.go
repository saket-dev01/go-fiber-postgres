package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct{
	Host		string	
	Port		string
	Password	string
	User 		string
	DBName		string
	SSLMode		string	
}

func NewConnection(config *Config)(*gorm.DB, error){
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	) // dsn is the string that will be needed by gorm to connect to the Database
	db, err:=gorm.Open(postgres.Open(dsn), &gorm.Config{}) // the .Open method takes the dsn string as input to connect to the DB

	return db, err
}