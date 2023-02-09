package data

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Form struct {
	_Id       primitive.ObjectID `bson:"_id"`
	Floor     string             `bson:"floor, omitempty"`
	Plinth    string             //плинтус
	Door      string
	Toilet    string
	Socket    string
	Plumb     string //сантехника
	Paint     string //краска
	Wallpaper string
	Tile      string //керамограит
	Cabinet   string //тумба под раковину
}

type FormModel struct {
	DB *mongo.Client
}

func (m FormModel) Insert(form *Form) error {
	collection := m.DB.Database("renovation").Collection("form")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := collection.InsertOne(ctx, form)
	if err != nil {
		return err
	}
	return nil
}

func (m FormModel) Get(id primitive.ObjectID) (*Form, error) {
	var form Form
	collection := m.DB.Database("test").Collection("form")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := collection.FindOne(ctx, id).Decode(&form)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, errors.New("record not found")
		default:
			return nil, err
		}
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
	result, err := collection.DeleteOne(ctx, form._Id)

	if err != nil || result.DeletedCount == 0 {
		switch {
		case result.DeletedCount == 0:
			return errors.New("record not found")
		default:
			return err
		}
	}
	return nil
}
