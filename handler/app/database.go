package app

import (
	"chera_khube/internal/database"
	"time"
)

func (a *application) openDatabaseConnection() error {
	db, err := database.Connect(a.config.Database)
	if err != nil {
		return err
	}
	a.db = db
	return nil
}

func (a *application) InitMemCache() database.MemoryCache {
	return database.NewMemoryCache(time.Hour*24, time.Hour*24)
}
