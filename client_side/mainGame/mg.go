package mainGame

import (
	"errors"
	"github.com/Jla3eP/tetris/both_sides_code"
	"github.com/Jla3eP/tetris/client_side/constants"
	"github.com/Jla3eP/tetris/client_side/field"
	"github.com/Jla3eP/tetris/client_side/render"
	"github.com/Jla3eP/tetris/client_side/requests"
	et "github.com/hajimehoshi/ebiten/v2"
	"math"
	"sync"
	"time"
)

type MainGame struct {
	render       *render.Render
	yourField    *field.Field
	enemiesField *field.Field
	figure       *field.Figure

	moveDownDefaultTime time.Duration
	moveDownShortTime   time.Duration
	moveDownTimer       time.Time
	moveTickerDuration  time.Duration
	moveTicker          *time.Ticker

	rotateTimer time.Time
	rotateSleep time.Duration

	movingDown  bool
	closeWindow bool

	scores      int
	enemyScores int
	status      int

	checkUpdatesTicker *time.Ticker
	gameStopper        sync.Once
}

func NewGame() *MainGame {
	et.SetMaxTPS(500)
	mg := &MainGame{}
	mg.yourField = field.NewField()
	mg.enemiesField = field.NewField()
	mg.render = createRenderObj(mg.yourField.GetSize())

	mg.moveDownDefaultTime = 220 * time.Millisecond
	mg.moveDownShortTime = 90 * time.Millisecond
	mg.moveDownTimer = time.Now()

	mg.moveTickerDuration = 100 * time.Millisecond
	mg.moveTicker = time.NewTicker(mg.moveTickerDuration)
	mg.rotateTimer = time.Now()
	mg.rotateSleep = 150 * time.Millisecond

	mg.movingDown = false
	mg.closeWindow = false

	mg.status = constants.StatusProcessing
	mg.checkUpdatesTicker = time.NewTicker(100 * time.Millisecond)
	return mg
}

func (mg *MainGame) Update() error {
	if mg.status == constants.StatusProcessing {
		requests.LogIn()
		mg.status = constants.StatusWaiting
		err := requests.FindGameRequest()
		if err != nil {
			return err
		}
		return nil
	} else if mg.status == constants.StatusWaiting {
		yourFigure, _, err := requests.GetGameInfoAndSendMyInfo(nil)
		if err != nil {
			return nil
		}
		mg.figure = yourFigure
		mg.status = constants.StatusPlaying
	} else if mg.status == constants.StatusPlaying {
		if mg.figure == nil || mg.figure.Fixed {
			reqInfo := &both_sides_code.FieldRequest{
				History: []both_sides_code.EnemyFigure{
					{
						EnemyFigureID:          int(mg.figure.GetID()),
						EnemyFigureRotateIndex: mg.figure.CurrentRotateIndex,
						EnemyFigureCoords:      mg.figure.CurrentCoords,
						EnemyFigureColor:       mg.figure.GetColor(),
					},
				},
			}
			yourFigure, enemyFigures, err := requests.GetGameInfoAndSendMyInfo(reqInfo)
			if err != nil {
				return err
			}
			mg.figure = yourFigure
			for _, v := range enemyFigures {
				mg.enemiesField.FixateFigure(v)
			}
			if mg.yourField.CheckCollision(mg.figure) {
				mg.status = constants.StatusWatching
			}
		}
	} else if mg.status == constants.StatusWatching {
		yourFigure, enemyFigures, err := requests.GetGameInfoAndSendMyInfo(nil)
		if err != nil {
			if err.Error() == "end" {
				mg.status = constants.StatusEnd
			}
		}
		mg.figure = yourFigure
		for _, v := range enemyFigures {
			mg.enemiesField.FixateFigure(v)
		}
	}
	if mg.status == constants.StatusPlaying {
		mg.processAll()
	} else if mg.status == constants.StatusWatching || mg.status == constants.StatusEnd {
		mg.gameStopper.Do(func() {

		})
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
		_ = mg.yourField.TryMoveFigure(mg.figure)
	default:
		break
	}

	if et.IsKeyPressed(et.KeyW) && time.Now().After(mg.rotateTimer.Add(mg.rotateSleep)) {
		_ = mg.yourField.TryRotateFigure(mg.figure)
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

	if !mg.figure.MoveDown(mg.yourField) {
		destroyedLines := mg.yourField.ClearField()
		enemiesDestroyedLines := mg.enemiesField.ClearField()
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

			mg.scores += mg.calcScoreBonus(destroyedLines)
			mg.enemyScores += mg.calcScoreBonus(enemiesDestroyedLines)
		}
		mg.figure.Fixed = true
	}
}

func (mg *MainGame) calcScoreBonus(destroyedLines int) int {
	return int(float64(destroyedLines*500) * (math.Pow(1.5, float64(destroyedLines-1))))
}

func (mg *MainGame) processAll() {
	mg.processInput()
	mg.processMoveDown()
}

func (mg *MainGame) Draw(screen *et.Image) {
	mg.render.RenderAll(screen, mg.yourField, mg.enemiesField, mg.figure, mg.scores, mg.enemyScores, mg.status)
}

func (mg *MainGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return render.GetX(), render.GetY()
}

func createRenderObj(fieldSize both_sides_code.Coords2) *render.Render {
	return &render.Render{FieldSize: fieldSize}
}
