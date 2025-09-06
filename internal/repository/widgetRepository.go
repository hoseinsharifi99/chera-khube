package repository

import "chera_khube/internal/model"

type WidgetRepository interface {
	Send(widget *model.DivarWidget, postToken, accessToken string, service string) error
	Delete(postToken, accessToken string, service string) error
}
