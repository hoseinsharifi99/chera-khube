package repository

import "chera_khube/internal/model"

type AppLogRepository interface {
	Insert(log *model.AppLog) error
}
