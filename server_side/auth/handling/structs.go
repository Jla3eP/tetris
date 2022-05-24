package handling

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	AuthInfo struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
	}

	SessionUpdateRequest struct {
		Nickname   string `json:"nickname"`
		SessionKey string `json:"session_key"`
	}

	sessionValues struct {
		userAgent  string
		username   string
		id         primitive.ObjectID
		createdAt  time.Time
		lastUpdate time.Time
	}
)
