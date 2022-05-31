package field

import (
	"errors"
	"github.com/Jla3eP/tetris/both_sides_code"
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

func (f *Field) CheckCollision(figure *Figure) bool {
	figureCoords := figure.GetRightCoords()
	for _, block := range figureCoords.Coords {
		if len(f.Field) <= block.Y || block.Y < 0 || len(f.Field[block.Y]) <= block.X || block.X < 0 || f.Field[block.Y][block.X].IsActive {
			return true
		}
	}
	return false
}

func (f *Field) GetSize() both_sides_code.Coords2 {
	crds := both_sides_code.Coords2{} // 🤤
	crds.Y = len(f.Field)
	if crds.Y == 0 {
		crds.X = 0
	} else {
		crds.X = len(f.Field[0])
	}
	return crds
}

func (f *Field) FixateFigure(figure *Figure) {
	figureCoords := figure.GetRightCoords()

	for _, block := range figureCoords.Coords {
		f.Field[block.Y][block.X].IsActive = true
		f.Field[block.Y][block.X].Color = figure.GetColor()
	}
}

func (f *Field) TryRotateFigure(figure *Figure) error {
	figure.rotate()
	if f.CheckCollision(figure) {
		figure.moveRight()
		if f.CheckCollision(figure) {
			figure.moveLeft()
			figure.moveLeft()
			if f.CheckCollision(figure) {
				figure.moveRight()
				figure.backRotate()
				return errors.New("can't rotate figure")
			}
		}
	}

	return nil
}

func (f *Field) TryMoveFigure(figure *Figure) error {
	var err error = nil
	err = f.tryMoveRight(figure)
	if err != nil {
		return err
	}
	err = f.tryMoveLeft(figure)
	if err != nil {
		return err
	}
	return nil
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

func (f *Field) tryMoveRight(figure *Figure) error {
	figure.Mutex.Lock()
	defer func() {
		figure.Mutex.Unlock()
	}()
	if et.IsKeyPressed(et.KeyD) {
		figure.moveRight()
		if f.CheckCollision(figure) {
			figure.moveLeft()
			return errors.New("can't move right")
		}
	}
	return nil
}

func (f *Field) tryMoveLeft(figure *Figure) error {
	figure.Mutex.Lock()
	defer func() {
		figure.Mutex.Unlock()
	}()
	if et.IsKeyPressed(et.KeyA) {
		figure.moveLeft()
		if f.CheckCollision(figure) {
			figure.moveRight()
			return errors.New("can't move left")
		}
	}
	return nil
}
