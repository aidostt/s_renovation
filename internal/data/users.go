package data

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"s_renovation.net/validator"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	Id        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	Name      string             `bson:"name"`
	Surname   string             `bson:"surname"`
	Phone     string             `bson:"phone"`
	Email     string             `bson:"email"`
	Password  password           `bson:"password"`
	Role      int                `bson:"role"`
}

type password struct {
	plaintext *string `json:"plaintext,omitempty"`
	hash      []byte  `json:"hash"`
}

type UserModel struct {
	DB *mongo.Client
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	ValidateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (m UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	collection := m.DB.Database("renovation").Collection("users")
	_, err := collection.InsertOne(ctx, user)
	fmt.Printf("this is user pass from insert %v", user.Password)
	if err != nil {
		//add unique constraint on email
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	collection := m.DB.Database("renovation").Collection("users")
	filter := bson.M{"email": email}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, errors.New("record not found")
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) GelAll() ([]User, error) {
	var users []User
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	collection := m.DB.Database("renovation").Collection("users")
	result, err := collection.Find(ctx, bson.M{}, nil)
	if err != nil {
		fmt.Println("first case")
		return nil, err
	}
	if err = result.All(ctx, &users); err != nil {
		fmt.Println("second case")
		return nil, err
	}
	return users, nil
}

func (m UserModel) Update() {
	return
}

func (m UserModel) Delete(user User) error {
	collection := m.DB.Database("renovation").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	result, err := collection.DeleteOne(ctx, user.Id)

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
