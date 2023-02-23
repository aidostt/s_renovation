package data

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Order struct {
	Id         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Phone      string             `bson:"phone"`
	Pack       string             `bson:"pack"`
	Additional bool               `bson:"additional"`
	Details    string             `bson:"details"`
}

type OrderModel struct {
	DB *mongo.Client
}

func (m OrderModel) Insert(form *Order) error {
	collection := m.DB.Database("renovation").Collection("order")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	form.Id = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, form)
	if err != nil {
		//server error response
		//db insert failed
		return err
	}
	return nil
}

func (m OrderModel) Get(id primitive.ObjectID) (*Form, error) {
	var form Form
	collection := m.DB.Database("renovation").Collection("order")
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

func (m OrderModel) GelAll() ([]Order, error) {
	var orders []Order
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	collection := m.DB.Database("renovation").Collection("order")
	result, err := collection.Find(ctx, bson.M{}, nil)
	if err != nil {
		fmt.Println("first case")
		return nil, err
	}
	if err = result.All(ctx, &orders); err != nil {
		fmt.Println("second case")
		return nil, err
	}
	return orders, nil
}

func (m OrderModel) Update(form *Order) error {
	return nil
}

func (m OrderModel) Delete(form *Order) error {
	collection := m.DB.Database("renovation").Collection("order")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	result, err := collection.DeleteOne(ctx, form.Id)

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
