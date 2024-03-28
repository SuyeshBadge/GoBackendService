package database

import (
	"backendService/internals/setup/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type dbConfig struct {
	dbType     string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
}

func (c *dbConfig) mysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.dbUser, c.dbPassword, c.dbHost, c.dbPort, c.dbName)
}

func (c *dbConfig) postgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.dbHost, c.dbUser, c.dbPassword, c.dbName, c.dbPort)
}

var Db *gorm.DB
var db *gorm.DB

func Connect(config *dbConfig) error {
	var err error

	switch config.dbType {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.mysqlDSN()), &gorm.Config{})

	case "postgres":
		db, err = gorm.Open(postgres.Open(config.postgresDSN()), &gorm.Config{})

	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.dbName), &gorm.Config{})

	default:
		return fmt.Errorf("unsupported database type: %s", config.dbType)
	}

	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return err
	}

	Db = db
	log.Printf("Connected to database :: %s", Db.Name())
	return nil
}

func InitializeDataBase(databaseType string) error {
	log.Println("Initializing database...")
	config := dbConfig{
		dbType:     databaseType,
		dbHost:     config.Config.Database.Host,
		dbPort:     config.Config.Database.Port,
		dbUser:     config.Config.Database.Username,
		dbPassword: config.Config.Database.Password,
		dbName:     config.Config.Database.Database,
	}

	if err := Connect(&config); err != nil {
		log.Fatal("Unable to connect database:", config.dbHost, ":", config.dbPort, "/", config.dbName)
		return err
	}
	log.Println("Database initialization completed.")
	return nil
}
