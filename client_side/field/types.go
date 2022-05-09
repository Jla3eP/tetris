package field

import (
	"sync"
)

type FieldPart struct {
	IsActive bool
	Color    int
}

type Field struct {
	Field [][]FieldPart
	rwm   *sync.RWMutex
}

type Coords2 struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type possibleStatus struct {
	Coords []Coords2 `json:"vec2"`
}

type FigureInterface interface {
	rotate()
	backRotate()
	GetCoords() *Coords2
	GetCurrentStatus() *possibleStatus
	GetColor() int
	GetRightCoords() *possibleStatus
	MoveRight()
	MoveLeft()
}

type Figure struct {
	PossibleStatuses   []possibleStatus `json:"possible_statuses"`
	Color              int
	CurrentRotateIndex int
	CurrentCoords      Coords2
}

type figuresConfig struct {
	Figures []Figure `json:"figures"`
}
