package data

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Form struct {
	id        int64
	floor     string
	plinth    string //плинтус
	door      string
	toilet    string
	socket    string
	plumb     string //сантехника
	paint     string //краска
	wallpaper string
	tile      string //керамограит
	cabinet   string //тумба под раковину
}

type FormModel struct {
	DB *mongo.Client
}

func (m FormModel) Insert(form *Form) (int, error) {
	collection := m.DB.Database("test").Collection("form")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	result, err := collection.InsertOne(ctx, form)
	if err != nil {
		return 0, err
	}
	return result.InsertedID.(int), nil
}

func (m FormModel) Get(id int64) (*Form, error) {
	var form Form
	collection := m.DB.Database("test").Collection("form")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := collection.FindOne(ctx, id).Decode(&form)
	if err != nil {
		switch {
		case errors.As(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		}
		}
		return nil, err
	}
	return &form, nil
}

func (m FormModel) Update(form *Form) error {
	return nil
}

func (m FormModel) Delete(form *Form) error {
	collection := m.DB.Database("test").Collection("form")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	result, err := collection.DeleteOne(ctx, form.id)
	if err != nil {
		return err
	}
	return nil
}
