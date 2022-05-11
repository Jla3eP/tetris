package mainGame

import (
	"github.com/Jla3eP/tetris/client_side/field"
	"github.com/Jla3eP/tetris/client_side/render"
	et "github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

type MainGame struct {
	render *render.Render
	field  *field.Field
	figure *field.Figure

	defaultMoveDownTicker *time.Ticker
	shortMoveDownTicker   *time.Ticker
	moveTicker            *time.Ticker
	rotateTicker          *time.Ticker

	moveDown bool
}

func NewGame() *MainGame {
	rand.Seed(time.Now().UnixNano())
	et.SetMaxTPS(60)
	mg := &MainGame{}
	mg.field = field.NewField()
	mg.render = createRenderObj(mg.field.GetSize())

	mg.defaultMoveDownTicker = time.NewTicker(350 * time.Millisecond)
	mg.shortMoveDownTicker = time.NewTicker(100 * time.Millisecond)
	mg.moveTicker = time.NewTicker(100 * time.Millisecond)
	mg.rotateTicker = time.NewTicker(250 * time.Millisecond)

	mg.moveDown = false
	return mg
}

func (mg *MainGame) Update() error {
	//<-mg.tickTicker.C
	if mg.figure == nil || mg.figure.Fixed {
		mg.figure = field.GetRandomFigure()
	}
	mg.processAll()

	return nil
}

func (mg *MainGame) processInput() {
	if !mg.moveDown && et.IsKeyPressed(et.KeyS) {
		mg.moveDown = true
	}
	if mg.moveDown && !et.IsKeyPressed(et.KeyS) {
		mg.moveDown = false
	}

	select {
	case <-mg.moveTicker.C:
		_ = mg.field.TryMoveFigure(mg.figure)
	default:
		break
	}

	if et.IsKeyPressed(et.KeyW) {
		select {
		case <-mg.rotateTicker.C:
			_ = mg.field.TryRotateFigure(mg.figure)
		default:
			break
		}
	}
}

func (mg *MainGame) processMoveDown() {
	if !mg.moveDown {
		select {
		case <-mg.defaultMoveDownTicker.C:
			break
		default:
			return
		}
	} else {
		select {
		case <-mg.shortMoveDownTicker.C:
			break
		default:
			return
		}
	}

	if !mg.figure.MoveDown(mg.field) {
		mg.field.ClearField()
		mg.figure.Fixed = true
		return
	}
}

func (mg *MainGame) processAll() {
	mg.processInput()
	mg.processMoveDown()
}

func (mg *MainGame) Draw(screen *et.Image) {
	mg.render.RenderAll(screen, mg.field, mg.figure)
}

func (mg *MainGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return render.GetX(), render.GetY()
}

func createRenderObj(fieldSize field.Coords2) *render.Render {
	return &render.Render{FieldSize: fieldSize}
}
