package repository

import "chera_khube/internal/model"

type ConfigRepository interface {
	List(serviceName string) []model.Config
	GetByCodes(codes []string) []model.Config
	ListAsMap() map[string]string
}
