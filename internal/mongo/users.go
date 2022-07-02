package mongo

import (
	"context"
	"fmt"
	"gitlab.com/g6834/team41/auth/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Users struct {
	c *mongo.Collection
}

func NewUsers(login, password, address, db string, port int) (*Users, error) {
	var ctx = context.TODO()
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", login, password, address, port, db)
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	u := &Users{
		c: client.Database(db).Collection("users"),
	}

	return u, nil
}

func (u Users) ChangeToken(token, login string) error {
	_, err := u.c.UpdateOne(context.TODO(), bson.M{"login": login}, bson.M{"$set": bson.M{"token": token}})
	if err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}

	return nil
}

func (u Users) GetUser(login string) (*models.User, error) {
	var user models.User
	err := u.c.FindOne(context.TODO(), bson.M{"login": login}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (u Users) AddUser(user models.User) error {
	_, err := u.c.InsertOne(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}

func (u Users) DeleteUser(login string) error {
	_, err := u.c.DeleteOne(context.TODO(), bson.M{"login": login})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
