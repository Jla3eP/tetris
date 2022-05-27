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

	moveDownDefaultTime time.Duration
	moveDownShortTime   time.Duration
	moveDownTimer       time.Time
	moveTickerDuration  time.Duration
	moveTicker          *time.Ticker

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

	mg.moveDownDefaultTime = 220 * time.Millisecond
	mg.moveDownShortTime = 90 * time.Millisecond
	mg.moveDownTimer = time.Now()

	mg.moveTickerDuration = 100 * time.Millisecond
	mg.moveTicker = time.NewTicker(mg.moveTickerDuration)
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
		if !time.Now().After(mg.moveDownTimer.Add(mg.moveDownDefaultTime)) {
			return
		}
	} else {
		if !time.Now().After(mg.moveDownTimer.Add(mg.moveDownShortTime)) {
			return
		}
	}

	mg.moveDownTimer = time.Now()

	if !mg.figure.MoveDown(mg.field) {
		destroyedLines := mg.field.ClearField()
		if destroyedLines != 0 {
			mg.moveDownDefaultTime /= 100
			mg.moveDownDefaultTime *= 98

			mg.moveDownShortTime /= 100
			mg.moveDownShortTime *= 98

			mg.rotateSleep /= 100
			mg.rotateSleep *= 98

			mg.moveTickerDuration /= 100
			mg.moveTickerDuration *= 98
			mg.moveTicker.Reset(mg.moveTickerDuration)

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
