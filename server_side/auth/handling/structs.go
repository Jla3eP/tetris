package handling

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
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

	lastStatuses struct {
		mu          *sync.RWMutex
		player1JSON []byte
		player2JSON []byte
	}

	gameInfo struct {
		ID              int64 `json:"id"`
		player1nickname string
		player2nickname string
	}
)
