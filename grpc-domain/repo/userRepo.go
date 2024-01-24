package repository

import (
	"context"
	"fmt"
	"ngc-grpc/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(input *model.UserAll) (string, error)
	ReadByName(name string) (*model.UserAll, error)
	ReadAll() ([]*model.UserModel, error)
	ReadID(ID string) (*model.User, error)
	Delete(ID string) (*mongo.DeleteResult, error)
	Update(oldID string, input *model.UserModel) (*mongo.UpdateResult, error)
}

type mongoRepository struct {
	DB *mongo.Client
}

func NewMongoRepository(db *mongo.Client) *mongoRepository {
	return &mongoRepository{
		DB: db,
	}
}

func (t *mongoRepository) Create(input *model.UserAll) (string, error) {
	ctx := context.TODO()
	db := t.DB.Database("grpc-test")

	res, err := db.Collection("Users").InsertOne(ctx, input)
	if err != nil {
		return "", err
	}
	oid, _ := res.InsertedID.(primitive.ObjectID)

	return oid.Hex(), nil
}

func (t *mongoRepository) ReadByName(name string) (*model.UserAll, error) {
	ctx := context.TODO()

	filter := bson.D{primitive.E{Key: "name", Value: name}}

	collection := t.DB.Database("grpc-test").Collection("Users")

	var result model.UserAll

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, fmt.Errorf("user not found")
	}

	return &result, nil
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

func (t *mongoRepository) ReadID(ID string) (*model.User, error) {
	ctx := context.TODO()

	filter := bson.D{primitive.E{Key: "id", Value: ID}}

	collection := t.DB.Database("grpc-test").Collection("Users")

	var result model.User

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &result, nil
}

func (t *mongoRepository) Delete(ID string) (*mongo.DeleteResult, error) {
	ctx := context.TODO()

	filter := bson.D{primitive.E{Key: "id", Value: ID}}

	collection := t.DB.Database("grpc-test").Collection("Users")

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (t *mongoRepository) Update(oldID string, input *model.UserModel) (*mongo.UpdateResult, error) {
	ctx := context.TODO()

	collection := t.DB.Database("grpc-test").Collection("Users")

	filter := bson.D{primitive.E{Key: "id", Value: oldID}}

	update := bson.M{
		"$set": bson.M{
			"id":   input.ID,
			"nama": input.Name,
		}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
