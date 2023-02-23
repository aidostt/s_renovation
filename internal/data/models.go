package data

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	Form  FormModel
	User  UserModel
	Order OrderModel
}

func NewModels(db *mongo.Client) Models {
	return Models{
		Form:  FormModel{DB: db},
		User:  UserModel{DB: db},
		Order: OrderModel{DB: db},
	}
}
