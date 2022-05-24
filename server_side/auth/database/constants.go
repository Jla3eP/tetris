package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection
var ctx = context.Background()

const (
	dbHost = "localhost"
	dbPort = "27017"

	idAndNicknameToSaltFormat = "%s+%s" //id.InsertedID.(primitive.ObjectID).String(), user.NickName
)
