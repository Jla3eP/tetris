package field

import "sync"

func (f *Figure) GetCopy() *Figure {
	figure := &Figure{PossibleStatuses: make([]PossibleStatus, len(f.PossibleStatuses))}
	for i := range f.PossibleStatuses {
		figure.PossibleStatuses[i] = *f.PossibleStatuses[i].GetCopy()
	}
	return figure
}

func (p *PossibleStatus) GetCopy() *PossibleStatus {
	PS := &PossibleStatus{Coords: make([]Coords2, len(p.Coords))}
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

	Coords2 struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	PossibleStatus struct {
		Coords []Coords2 `json:"vec2"`
	}

	Figure struct {
		id                 int8
		PossibleStatuses   []PossibleStatus `json:"possible_statuses"`
		Color              int
		CurrentRotateIndex int
		CurrentCoords      Coords2
		Mutex              *sync.Mutex
		Fixed              bool
	}

	figuresConfig struct {
		Figures []Figure `json:"figures"`
	}
)
