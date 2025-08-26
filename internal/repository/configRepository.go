package repository

import "chera_khube/internal/model"

type ConfigRepository interface {
	List() []model.Config
	GetByCodes(codes []string) []model.Config
	ListAsMap() map[string]string
}
