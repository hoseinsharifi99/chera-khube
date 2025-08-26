package main

import (
	"chera_khube/handler/app"
	"chera_khube/internal/constant"
	"chera_khube/internal/database"
	"chera_khube/internal/model"
	"flag"
	"gorm.io/gorm"
)

func main() {
	path := flag.String("e", constant.DefaultEnvPath, "env file path")
	flag.Parse()
	config, err := app.SetupViper(*path)
	if err != nil {
		panic(err)
	}

	db, err := database.Connect(config.Database)
	if err != nil {
		panic(err)
	}
	migrate(db)
}

func migrate(db *gorm.DB) {
	db.Debug().AutoMigrate(model.Post{})
	db.Debug().AutoMigrate(model.UserPayment{})
	db.Debug().AutoMigrate(model.PricingLogic{})
	db.Debug().AutoMigrate(model.Adons{})
	db.Debug().AutoMigrate(model.User{})
}
