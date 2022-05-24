package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/Jla3eP/tetris/server_side/auth/hash"
	usr "github.com/Jla3eP/tetris/server_side/auth/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func UserExists(ctx context.Context, userName string) bool {
	filter := bson.M{"nickname": userName}
	resp := collection.FindOne(ctx, filter)
	if resp.Err() != nil {
		return false
	}
	return true
}

func CreateUser(ctx context.Context, user usr.User, clearPassword string) error {
	if UserExists(ctx, user.NickName) {
		return errors.New(fmt.Sprintf("user with nickname=\"%s\" exists", user.NickName))
	}

	userData := bson.M{
		"nickname":       user.NickName,
		"account_status": usr.StatusActive,
		"created_at":     time.Now(),
	}

	id, err := collection.InsertOne(ctx, userData)
	if err != nil {
		return err
	}

	filter := bson.D{{"nickname", user.NickName}}
	user.ID = id.InsertedID.(primitive.ObjectID)
	salt := infoToSalt(user)

	userData["hashed_password"] = hash.CreateSaltPasswordHash(salt, clearPassword)
	update := bson.D{
		{"$set", userData},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func VerifyPassword(ctx context.Context, user usr.User, clearPassword string) (bool, error) {
	filter := bson.D{{"nickname", user.NickName}}

	if user.HashedPassword == "" {
		err := collection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			return false, err
		}
	}

	salt := infoToSalt(user)
	hashPassword := hash.CreateSaltPasswordHash(salt, clearPassword)

	if hashPassword != user.HashedPassword {
		return false, errors.New("invalid password")
	}

	return true, nil
}

func GetIdByUsername(username string) (primitive.ObjectID, error) {
	user := usr.User{}
	filter := bson.D{{"nickname", username}}

	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return [12]byte{}, err
	}

	return user.ID, nil
}

func infoToSalt(usr usr.User) string {
	return fmt.Sprintf(idAndNicknameToSaltFormat, usr.ID.String(), usr.NickName)
}

func init() {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s/", dbHost, dbPort))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("tetris").Collection("user")
}
