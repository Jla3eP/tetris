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

type PossibleStatus struct {
	Coords []Coords2 `json:"vec2"`
}

func (p *PossibleStatus) GetCopy() *PossibleStatus {
	PS := &PossibleStatus{Coords: make([]Coords2, len(p.Coords))}
	copy(PS.Coords, p.Coords)
	return PS
}

type Figure struct {
	id                 int8
	PossibleStatuses   []PossibleStatus `json:"possible_statuses"`
	Color              int
	CurrentRotateIndex int
	CurrentCoords      Coords2
	Mutex              *sync.Mutex
	Fixed              bool
}

func (f *Figure) GetCopy() *Figure {
	figure := &Figure{PossibleStatuses: make([]PossibleStatus, len(f.PossibleStatuses))}
	for i := range f.PossibleStatuses {
		figure.PossibleStatuses[i] = *f.PossibleStatuses[i].GetCopy()
	}
	return figure
}

type figuresConfig struct {
	Figures []Figure `json:"figures"`
}
