package field

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var figures []Figure

func GetFigures() []Figure {
	return figures
}

func (f *Figure) MoveDown(field Field) {
	f.CurrentCoords.Y++

	if field.CheckCollision(FigureInterface(f)) {
		f.CurrentCoords.Y--
		field.FixateFigure(f)
	}
}

func (f *Figure) GetCoords() *Coords2 {
	return &f.CurrentCoords
}

func (f *Figure) GetCurrentStatus() *possibleStatus {
	return &f.PossibleStatuses[f.CurrentRotateIndex]
}

func (f *Figure) GetRightCoords() *possibleStatus {
	currentState := *f.GetCurrentStatus()
	coords := f.GetCoords()

	statusWithRightCoords := possibleStatus{currentState.Coords[:]}

	for i := range statusWithRightCoords.Coords {
		statusWithRightCoords.Coords[i].X += coords.X
		statusWithRightCoords.Coords[i].Y += coords.Y
	}

	return &statusWithRightCoords
}

func (f *Figure) GetColor() int {
	return f.Color
}

func (f *Figure) MoveRight() {
	f.CurrentCoords.X++
}

func (f *Figure) MoveLeft() {
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
	fmt.Println(config)

	figures = config.Figures
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
