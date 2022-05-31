package field

import (
	"github.com/Jla3eP/tetris/both_sides_code"
	"sync"
)

func (f *Figure) GetCopy() *Figure {
	figure := &Figure{PossibleStatuses: make([]PossibleStatus, len(f.PossibleStatuses))}
	for i := range f.PossibleStatuses {
		figure.PossibleStatuses[i] = *f.PossibleStatuses[i].GetCopy()
	}
	return figure
}

func (p *PossibleStatus) GetCopy() *PossibleStatus {
	PS := &PossibleStatus{Coords: make([]both_sides_code.Coords2, len(p.Coords))}
	copy(PS.Coords, p.Coords)
	return PS
}

type (
	FieldPart struct {
		IsActive bool
		Color    int
	}

	Field struct {
		Field [][]FieldPart
		rwm   *sync.RWMutex
	}

	PossibleStatus struct {
		Coords []both_sides_code.Coords2 `json:"vec2"`
	}

	Figure struct {
		id                 int8
		PossibleStatuses   []PossibleStatus `json:"possible_statuses"`
		Color              int
		CurrentRotateIndex int
		CurrentCoords      both_sides_code.Coords2
		Mutex              *sync.Mutex
		Fixed              bool
	}

	figuresConfig struct {
		Figures []Figure `json:"figures"`
	}
)
