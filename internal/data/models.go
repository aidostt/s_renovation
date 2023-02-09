package data

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Form FormModel
	User UserModel
}

func NewModels(db *mongo.Client) Models {
	return Models{
		Form: FormModel{DB: db},
		User: UserModel{DB: db},
	}
}
