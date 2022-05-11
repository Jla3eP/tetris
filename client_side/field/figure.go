package field

import (
	"encoding/json"
	"github.com/Jla3eP/tetris/client_side/constants"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sync"
)

var (
	figures         []Figure
	randFigureIndex int
	randColorIndex  int
	colors          []int
	mutex           *sync.Mutex
)

func GetFigures() []Figure {
	figuresCopy := make([]Figure, len(figures))
	copy(figuresCopy, figures)
	return figuresCopy
}

func GetRandomFigure() *Figure {
	figures = GetFigures()
	randFigureIndex = rand.Int() % len(figures)
	randColorIndex = rand.Int() % len(colors)
	figure := figures[randFigureIndex].GetCopy()
	figure.Color = colors[randColorIndex]
	figure.Mutex = mutex
	figure.CurrentCoords.X = 4
	figure.Fixed = false
	return figure
}

func (f *Figure) MoveDown(field *Field) bool {
	f.CurrentCoords.Y++

	if field.CheckCollision(f) {
		f.CurrentCoords.Y--
		field.FixateFigure(f)
		return false
	}
	return true
}

func (f *Figure) GetCurrentStatus() PossibleStatus {
	ps := PossibleStatus{}
	ps.Coords = f.PossibleStatuses[f.CurrentRotateIndex].Coords[:]
	return ps
}

func (f *Figure) GetRightCoords() *PossibleStatus {
	currentState := f.GetCurrentStatus()
	coords := f.CurrentCoords

	statusWithRightCoords := PossibleStatus{Coords: make([]Coords2, len(currentState.Coords))}
	copy(statusWithRightCoords.Coords, currentState.Coords)

	for i := range statusWithRightCoords.Coords {
		statusWithRightCoords.Coords[i].X += coords.X
		statusWithRightCoords.Coords[i].Y += coords.Y
	}
	return &statusWithRightCoords
}

func (f *Figure) GetColor() int {
	return f.Color
}

func (f *Figure) moveRight() {
	f.CurrentCoords.X++
}

func (f *Figure) moveLeft() {
	f.CurrentCoords.X--
}

func init() {
	file, err := os.Open("client_side/field/figures_config.json")
	defer file.Close()
	if err != nil {
		log.Fatalln(err)
		return
	}

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
		return
	}
	config := figuresConfig{}

	err = json.Unmarshal(buffer, &config)
	if err != nil {
		log.Fatalln(err.Error() + " module figure (init)")
		return
	}

	figures = config.Figures
	for i := range figures {
		figures[i].id = int8(i + 1)
	}
	colors = append(colors,
		constants.ColorBlue,
		constants.ColorGreen,
		constants.ColorOrange,
		constants.ColorRed,
		constants.ColorYellow)
	randFigureIndex = -1
	mutex = &sync.Mutex{}
}

func (f *Figure) rotate() {
	f.CurrentRotateIndex++
	if f.CurrentRotateIndex >= len(f.PossibleStatuses) {
		f.CurrentRotateIndex = 0
	}
}

func (f *Figure) backRotate() {
	f.CurrentRotateIndex--
	if f.CurrentRotateIndex <= -1 {
		f.CurrentRotateIndex = len(f.PossibleStatuses) - 1
	}
}
