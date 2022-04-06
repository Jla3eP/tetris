package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection
var ctx = context.Background()

const host = "localhost"
const port = "27017"

const idAndNicknameToSaltFormat = "%s+%s" //id.InsertedID.(primitive.ObjectID).String(), user.NickName
