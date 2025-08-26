package response

import "chera_khube/internal/model"

type UserBalancePayment struct {
	Balance      int                  `json:"balance"`
	UserPayments []*model.UserPayment `json:"user_payments"`
}
