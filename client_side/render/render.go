package render

import (
	"encoding/json"
	"fmt"
	"github.com/Jla3eP/tetris/client_side/constants"
	"github.com/Jla3eP/tetris/client_side/field"
	et "github.com/hajimehoshi/ebiten/v2"
	drw "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	fnt "golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	windowX int
	windowY int

	textures map[int]*et.Image

	textureWidth  int32
	textureHeight int32

	bonusPxWindowWidth        int32
	bonusPxWindowHeight       int32
	bonusPercentsWindowWidth  int32
	bonusPercentsWindowHeight int32

	fieldBiasX           int32
	fieldBiasY           int32
	enemyFieldBiasX      int32
	pixelsToEnemiesField int32

	font                 fnt.Face
	enemyFieldBiasSetter sync.Once

	waitingAnimationTicker = time.NewTicker(250 * time.Millisecond)
	postfixWaitingString   = "."
	waitingStringMu        = &sync.RWMutex{}
	waitingStringUpdater   = sync.Once{}
	animStopper            = make(chan struct{})
)

const pathToTexturesFormat = "render/textures/%s/texture_sq_%s.png"

func GetX() int {
	return windowX
}

func GetY() int {
	return windowY
}

func updateWaitingPostfix() {
loop:
	for {
		select {
		case <-animStopper:
			waitingAnimationTicker.Stop()
			break loop
		case <-waitingAnimationTicker.C:
			waitingStringMu.Lock()
			postfixWaitingString += "."
			if len(postfixWaitingString) > 4 {
				postfixWaitingString = "."
			}
			waitingStringMu.Unlock()
		}
	}
}

func (r *Render) RenderAll(screen *et.Image, field *field.Field, enemyField *field.Field, figure *field.Figure, scores, enemyScores int, status int) {
	enemyFieldBiasSetter.Do(func() {
		r.setEnemyFieldBias()
	})
	if status == constants.StatusWaiting {
		r.renderWaiting(screen, field)
	}
	r.renderFieldBackground(screen, r.FieldSize.X, r.FieldSize.Y, fieldBiasX, fieldBiasY)
	r.renderEnemyFieldBackground(screen)
	r.renderField(screen, field, fieldBiasX, fieldBiasY)
	r.renderEnemiesField(screen, enemyField)
	if status == constants.StatusPlaying {
		r.renderFigure(screen, figure, fieldBiasX, fieldBiasY)
		r.renderInfo(screen, scores, field)
		r.renderEnemyInfo(screen, enemyScores, enemyField)
	} else if status == constants.StatusWatching {
		r.renderEndgameInfo(screen, scores, field)
	}
}

func (r *Render) renderWaiting(screen *et.Image, f *field.Field) {
	waitingStringMu.RLock()
	text.Draw(screen, "Waiting"+postfixWaitingString, font, int(fieldBiasX+(textureWidth+1)*int32(f.GetSize().X)),
		int(fieldBiasY+textureHeight/2), color.RGBA{R: 255, G: 255, B: 255, A: 255})
	waitingStringMu.RUnlock()
}

func (r *Render) renderEndgameInfo(screen *et.Image, scores int, f *field.Field) { //TODO
	text.Draw(screen, "Game over\nYour scores:"+strconv.Itoa(scores)+"\nPress ESC to exit", font, int(fieldBiasX+(textureWidth+1)*int32(f.GetSize().X)),
		int(fieldBiasY), color.RGBA{R: 255, G: 255, B: 255, A: 255})
}

func (r *Render) renderEnemyFieldBackground(screen *et.Image) {
	r.renderFieldBackground(screen, r.FieldSize.X, r.FieldSize.Y, enemyFieldBiasX, fieldBiasY)
}

func (r *Render) renderFieldBackground(screen *et.Image, x, y int, biasX, biasY int32) {
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			r.renderFieldElement(
				screen,
				biasX+int32(i)*textureWidth,
				biasY+int32(j)*textureHeight,
				constants.ColorBlack)
		}
	}

}

func (r *Render) renderFieldElement(screen *et.Image, x, y int32, color int) {
	options := &et.DrawImageOptions{}
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(textures[color], options)
}

func (r *Render) renderEnemiesField(screen *et.Image, enemyField *field.Field) {
	r.renderField(screen,
		enemyField,
		fieldBiasX+textureWidth*int32(r.FieldSize.X)+pixelsToEnemiesField,
		fieldBiasY,
	)
}

func (r *Render) renderField(screen *et.Image, f *field.Field, biasX, biasY int32) {
	if f == nil {
		return
	}
	for y := range f.Field {
		for x := range f.Field[y] {
			if f.Field[y][x].IsActive {
				r.renderFieldElement(
					screen,
					biasX+int32(x)*textureWidth,
					biasY+int32(y)*textureHeight,
					f.Field[y][x].Color)
			}
		}
	}
}

func (r *Render) renderFigure(screen *et.Image, f *field.Figure, biasX, biasY int32) {
	stateWithRightCoords := f.GetRightCoords()

	for i := range stateWithRightCoords.Coords {
		r.renderFieldElement(screen,
			biasX+int32(stateWithRightCoords.Coords[i].X)*textureWidth,
			biasY+int32(stateWithRightCoords.Coords[i].Y)*textureHeight,
			f.GetColor())
	}
}

func (r *Render) renderInfo(screen *et.Image, scores int, f *field.Field) {
	text.Draw(screen, strconv.Itoa(scores), font, int(fieldBiasX+(textureWidth+1)*int32(f.GetSize().X)),
		int(fieldBiasY), color.RGBA{R: 255, G: 255, B: 255, A: 255})
}

func (r *Render) renderEnemyInfo(screen *et.Image, scores int, f *field.Field) {
	text.Draw(screen, strconv.Itoa(scores), font, int(enemyFieldBiasX+(textureWidth+1)*int32(f.GetSize().X)),
		int(fieldBiasY), color.RGBA{R: 255, G: 255, B: 255, A: 255})
}

func addTextureToMap(key int, TexturePackName, name string) {
	var err error
	textures[key], _, err = drw.NewImageFromFile(fmt.Sprintf(pathToTexturesFormat, TexturePackName, name))
	if err != nil {
		log.Fatalln(err)
	}
}

func (r *Render) setEnemyFieldBias() {
	enemyFieldBiasX = fieldBiasX + textureWidth*int32(r.FieldSize.X) + pixelsToEnemiesField
}

func setConfigAndUploadTextures() {
	file, err := os.Open("./render/render_config.json")
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

	config := renderConfig{}
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		log.Fatalln(err)
		return
	}
	textureHeight = config.TextureHeight
	textureWidth = config.TextureWidth
	bonusPxWindowHeight = config.BonusPxWindowHeight
	bonusPxWindowWidth = config.BonusPxWindowWidth
	bonusPercentsWindowWidth = config.BonusPercentsWindowWidth
	bonusPercentsWindowHeight = config.BonusPercentsWindowHeight
	fieldBiasX = config.FieldBiasX
	fieldBiasY = config.FieldBiasY
	pixelsToEnemiesField = config.PixelsToEnemiesField

	windowX = int((int32(field.FieldX)*textureWidth)*(1+bonusPercentsWindowWidth/100))*2 + int(bonusPxWindowWidth+fieldBiasX)
	windowY = int((int32(field.FieldY)*textureHeight)*(1+bonusPercentsWindowHeight/100)) + int(bonusPxWindowHeight+fieldBiasY)

	textures = make(map[int]*et.Image)

	addTextureToMap(constants.ColorBlack, config.TexturePackName, "black")
	addTextureToMap(constants.ColorBlue, config.TexturePackName, "blue")
	addTextureToMap(constants.ColorGreen, config.TexturePackName, "green")
	addTextureToMap(constants.ColorOrange, config.TexturePackName, "orange")
	addTextureToMap(constants.ColorRed, config.TexturePackName, "red")
	addTextureToMap(constants.ColorYellow, config.TexturePackName, "yellow")

	font = basicfont.Face7x13
}

func init() {
	waitingStringUpdater.Do(func() {
		go updateWaitingPostfix()
	})
	setConfigAndUploadTextures()
}
