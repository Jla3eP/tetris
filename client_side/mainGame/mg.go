package mainGame

import (
	"errors"
	"github.com/Jla3eP/tetris/client_side/field"
	"github.com/Jla3eP/tetris/client_side/render"
	et "github.com/hajimehoshi/ebiten/v2"
	"math"
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

	rotateTimer time.Time
	rotateSleep time.Duration

	movingDown  bool
	endGame     bool
	closeWindow bool

	scores int
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
	mg.rotateTimer = time.Now()
	mg.rotateSleep = 150 * time.Millisecond

	mg.movingDown = false
	mg.closeWindow = false
	return mg
}

func (mg *MainGame) Update() error {
	if mg.figure == nil || mg.figure.Fixed {
		mg.figure = field.GetRandomFigure()
		if mg.field.CheckCollision(mg.figure) {
			mg.endGame = true
		}
	}
	if !mg.endGame {
		mg.processAll()
	} else {
		mg.checkESC()
	}

	if mg.closeWindow {
		return errors.New("it's okay. Game over")
	}

	return nil
}

func (mg *MainGame) checkESC() {
	if et.IsKeyPressed(et.KeyEscape) {
		mg.closeWindow = true
	}
}

func (mg *MainGame) processInput() {
	if !mg.movingDown && et.IsKeyPressed(et.KeyS) {
		mg.movingDown = true
	}
	if mg.movingDown && !et.IsKeyPressed(et.KeyS) {
		mg.movingDown = false
	}

	select {
	case <-mg.moveTicker.C:
		_ = mg.field.TryMoveFigure(mg.figure)
	default:
		break
	}

	if et.IsKeyPressed(et.KeyW) && time.Now().After(mg.rotateTimer.Add(mg.rotateSleep)) {
		_ = mg.field.TryRotateFigure(mg.figure)
		mg.rotateTimer = time.Now()
	}
}

func (mg *MainGame) processMoveDown() {
	if !mg.movingDown {
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
		destroyedLines := mg.field.ClearField()
		if destroyedLines != 0 {
			mg.scores += int(float64(destroyedLines*500) * (math.Pow(1.5, float64(destroyedLines-1))))
		}
		mg.figure.Fixed = true
		return
	}
}

func (mg *MainGame) processAll() {
	mg.processInput()
	mg.processMoveDown()
}

func (mg *MainGame) Draw(screen *et.Image) {
	mg.render.RenderAll(screen, mg.field, mg.figure, mg.scores, mg.endGame)
}

func (mg *MainGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return render.GetX(), render.GetY()
}

func createRenderObj(fieldSize field.Coords2) *render.Render {
	return &render.Render{FieldSize: fieldSize}
}
