package db

import (
	"dankmuzikk/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var instance *gorm.DB = nil

func Connector() (*gorm.DB, error) {
	return getDBConnector(config.Env().DB.Name)
}

func getDBConnector(dbName string) (*gorm.DB, error) {
	if instance != nil {
		return instance, nil
	}

	createDBDsn := fmt.Sprintf("%s:%s@tcp(%s)/", config.Env().DB.Username, config.Env().DB.Password, config.Env().DB.Host)
	database, err := gorm.Open(mysql.Open(createDBDsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = database.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + ";").Debug().Error
	if err != nil {
		return nil, err
	}
	instance, err := gorm.Open(mysql.Open(
		fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True&loc=Local",
			config.Env().DB.Username,
			config.Env().DB.Password,
			config.Env().DB.Host,
			dbName,
		),
	), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return instance, nil
}
