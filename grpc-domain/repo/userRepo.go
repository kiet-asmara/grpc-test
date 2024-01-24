package repository

import (
	"context"
	"ngc-grpc/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(input *model.UserModel) (string, error)
	ReadAll() ([]*model.UserModel, error)
}

type mongoRepository struct {
	DB *mongo.Client
}

func NewMongoRepository(db *mongo.Client) *mongoRepository {
	return &mongoRepository{
		DB: db,
	}
}

func (t *mongoRepository) Create(input *model.UserModel) (string, error) {
	ctx := context.TODO()
	db := t.DB.Database("grpc-test")

	res, err := db.Collection("Users").InsertOne(ctx, input)
	if err != nil {
		return "", err
	}
	oid, _ := res.InsertedID.(primitive.ObjectID)

	return oid.Hex(), nil
}

func (t *mongoRepository) ReadAll() ([]*model.UserModel, error) {
	ctx := context.TODO()

	collection := t.DB.Database("grpc-test").Collection("Users")

	var results []*model.UserModel

	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var user *model.UserModel
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}

		results = append(results, user)
	}

	return results, nil
}
