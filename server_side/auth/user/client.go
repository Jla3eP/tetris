package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	StatusActive  = 1
	StatusBanned  = 2
	StatusDeleted = 3
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt      time.Time          `bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty"`
	NickName       string             `bson:"nickname"`
	HashedPassword string             `bson:"hashed_password,omitempty"`
	AccountStatus  int                `bson:"account_status,omitempty"`
}
