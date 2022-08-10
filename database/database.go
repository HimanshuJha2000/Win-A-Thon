package database

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"win-a-thon/config"
	"win-a-thon/models"
)

var DB *gorm.DB

type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func BuildDBConfig(db config.Database) *DBConfig {
	dbconfig := DBConfig{
		db.Host,
		db.Port,
		db.Username,
		db.Name,
		db.Password,
	}

	return &dbconfig
}

func DBURL(dbConfig *DBConfig) string {
	fmt.Println("db config is: ", dbConfig.Port)
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}

func GetDatabase() (*gorm.DB, error) {
	appConfig := config.GetConfig()
	if _, err := toml.DecodeFile("config/env.default.toml", &appConfig); err != nil {
		fmt.Println(err)
		return nil, err
	}
	db, err := gorm.Open(mysql.Open(DBURL(BuildDBConfig(appConfig.Database))), &gorm.Config{})
	//db.DB().SetMaxIdleConns(databaseConfig.MaxIdleConnections)
	//db.DB().SetMaxOpenConns(databaseConfig.MaxOpenConnections)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Hackathon{}, &models.Participant{}, &models.Notification{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
