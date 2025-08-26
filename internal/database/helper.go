package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Config struct {
	Host              string `mapstructure:"host"`
	Port              string `mapstructure:"port"`
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	DBName            string `mapstructure:"db_name"`
	MaxOpenConnection int    `mapstructure:"max_idle_connection"`
	MaxIdleConnection int    `mapstructure:"max_open_connection"`
	MaxLifeTime       int    `mapstructure:"max_life_time"`
}

func Connect(config Config) (db *gorm.DB, err error) {
	//open database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
		config.Host,
		config.Username,
		config.Password,
		config.DBName,
		config.Port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	d, _ := db.DB()
	d.SetMaxOpenConns(config.MaxOpenConnection)
	d.SetMaxIdleConns(config.MaxIdleConnection)
	d.SetConnMaxLifetime(time.Second * time.Duration(config.MaxLifeTime))

	return db, err
}

func Close(db *gorm.DB) error {
	d, _ := db.DB()
	return d.Close()
}
