package field

import (
	"errors"
	et "github.com/hajimehoshi/ebiten/v2"
)

var (
	FieldX = uint(10)
	FieldY = uint(20)
)

func NewField() *Field {
	fld := &Field{}
	fld.Field = make([][]FieldPart, FieldY)
	for i := 0; i < len(fld.Field); i++ {
		fld.Field[i] = make([]FieldPart, 10)
	}

	return fld
}

func (f *Field) CheckCollision(figure FigureInterface) bool {
	figureCoords := figure.GetRightCoords()

	for _, block := range figureCoords.Coords {
		if len(f.Field) <= block.Y || len(f.Field[block.Y]) <= block.X || f.Field[block.Y][block.X].IsActive {
			return true
		}
	}
	return false
}

func (f *Field) GetSize() Coords2 {
	crds := Coords2{} // 🤤
	crds.Y = len(f.Field)
	if crds.Y == 0 {
		crds.X = 0
	} else {
		crds.X = len(f.Field[0])
	}
	return crds
}

func (f *Field) FixateFigure(figure FigureInterface) {
	figureCoords := figure.GetRightCoords()

	f.rwm.Lock()
	for _, block := range figureCoords.Coords {
		f.Field[block.Y][block.X].IsActive = true
		f.Field[block.Y][block.X].Color = figure.GetColor()
	}
	f.rwm.Unlock()
}

func (f *Field) TryRotateFigure(figure FigureInterface) error {
	figure.rotate()
	if f.CheckCollision(figure) {
		figure.MoveRight()
		if f.CheckCollision(figure) {
			figure.MoveLeft()
			figure.MoveLeft()
			if f.CheckCollision(figure) {
				figure.MoveRight()
				figure.backRotate()
				return errors.New("can't rotate figure")
			}
		}
	}
	return nil
}

func (f *Field) TryMoveFigure(figure FigureInterface) error {
	var err error = nil
	if et.IsKeyPressed(et.KeyD) {
		figure.MoveRight()
		if f.CheckCollision(figure) {
			figure.MoveLeft()
			err = errors.New("can't move right")
		}
	} else if et.IsKeyPressed(et.KeyA) {
		figure.MoveLeft()
		if f.CheckCollision(figure) {
			figure.MoveRight()
			err = errors.New("can't move left")
		}
	}

	return err
}

func (f *Field) ClearField() int {
	destroyedLines := 0
	for i := len(f.Field) - 1; i >= 0; i-- {
		goNext := false
		for j := 0; j < len(f.Field[i]); j++ {
			if !f.Field[i][j].IsActive {
				goNext = true
				break
			}
		}

		if goNext {
			continue
		}
		destroyedLines++
		f.clearLine(i)
		for t := i - 1; t >= 0; t-- {
			f.moveLineDown(t)
		}
		f.clearLine(0)
		i++
	}
	return destroyedLines
}

func (f *Field) moveLineDown(lineIndex int) {
	for i := 0; i < len(f.Field[lineIndex]); i++ {
		f.Field[lineIndex+1][i] = f.Field[lineIndex][i]
	}
}

func (f *Field) clearLine(lineIndex int) {
	for i := 0; i < len(f.Field[lineIndex]); i++ {
		f.Field[lineIndex][i].IsActive = false
	}
}
