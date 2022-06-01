package handling

import (
	"github.com/Jla3eP/tetris/both_sides_code"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func getNewGameInfo() gameInfo {
	atomic.AddInt64(&currentGameID, 1)
	return gameInfo{
		ID:                    currentGameID,
		players:               make([]string, 0, 2),
		status:                gameStatusWaiting,
		playerActive:          make([]bool, 2, 2),
		playerWatching:        make([]bool, 2, 2),
		playersFiguresIndexes: make([]int, 2, 2),
		lastStatuses: lastStatuses{
			mu:                  &sync.RWMutex{},
			playersLastStatuses: make([][]*both_sides_code.FieldRequest, 2),
		},
	}
}

func generateFigures(gi gameInfo) gameInfo {
	rand.Seed(time.Now().UnixNano())
	figuresProbabilities := make([]float64, figuresCount)
	for i := range figuresProbabilities {
		figuresProbabilities[i] = 1.0 / float64(figuresCount)
	}

	lastColor := -1
	for i := 0; i < 100; i++ {
		randFigureIndex := -1
		randFigureIndex = rand.Int() % 100000

		sum := 0
		for k, v := range figuresProbabilities {
			sum += int(v * 100000)
			if sum > randFigureIndex {
				randFigureIndex = k
				figuresProbabilities[k] /= 2.0
				for kk := range figuresProbabilities {
					if kk != k {
						figuresProbabilities[kk] += figuresProbabilities[k] / float64(figuresCount-1)
					}
				}
				break
			}
		}
		if randFigureIndex == -1 {
			gi.figures = append(gi.figures, figuresCount)
		} else {
			gi.figures = append(gi.figures, randFigureIndex)
		}

		randC := rand.Int() % figureColorsCount
		for lastColor == randC {
			randC = rand.Int() % figureColorsCount
		}
		lastColor = randC
		gi.figuresColors = append(gi.figuresColors, randC)
	}

	return gi
}

type (
	sessionValues struct {
		userAgent  string
		username   string
		id         primitive.ObjectID
		createdAt  time.Time
		lastUpdate time.Time
	}

	lastStatuses struct {
		mu                  *sync.RWMutex
		playersLastStatuses [][]*both_sides_code.FieldRequest
	}

	gameInfo struct {
		ID                    int64 `json:"currentGameID"`
		players               []string
		playerActive          []bool
		playerWatching        []bool
		status                int
		figures               []int
		figuresColors         []int
		playersFiguresIndexes []int
		lastStatuses
	}
)
