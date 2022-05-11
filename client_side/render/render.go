package render

import (
	"encoding/json"
	"fmt"
	"github.com/Jla3eP/tetris/client_side/constants"
	"github.com/Jla3eP/tetris/client_side/field"
	et "github.com/hajimehoshi/ebiten/v2"
	drw "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
)

var (
	TargetFps = int32(60)
	DrawFps   = true

	windowX = 0
	windowY = 0

	textures map[int]*et.Image

	textureWidth  = int32(25)
	textureHeight = int32(25)

	bonusPxWindowWidth        = int32(20)
	bonusPxWindowHeight       = int32(0)
	bonusPercentsWindowWidth  = int32(5)
	bonusPercentsWindowHeight = int32(5)

	fieldBiasX = uint32(50)
	fieldBiasY = uint32(50)
	endGame    = false
)

const pathToTexturesFormat = "client_side/render/textures/%s/texture_sq_%s.png"

func GetX() int {
	return windowX
}

func GetY() int {
	return windowY
}

func (r *Render) RenderAll(screen *et.Image, field *field.Field, figure *field.Figure /*rdCh <-chan *RenderData, endCh <-chan struct{}*/) { //TODO
	r.renderFieldBackground(screen, r.FieldSize.X, r.FieldSize.Y)
	r.renderField(screen, field)
	r.renderFigure(screen, figure)
}

func (r *Render) renderEndgameInfo(screen *et.Image) { //TODO
}

func (r *Render) renderFieldBackground(screen *et.Image, x, y int) {
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			r.renderFieldElement(
				screen,
				int32(fieldBiasX)+int32(i)*textureWidth,
				int32(fieldBiasY)+int32(j)*textureHeight,
				constants.ColorBlack)
		}
	}

}

func (r *Render) renderFieldElement(screen *et.Image, x, y int32, color int) {
	options := &et.DrawImageOptions{}
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(textures[color], options)
}

func (r *Render) renderField(screen *et.Image, f *field.Field) {
	if f == nil {
		return
	}
	for y := range f.Field {
		for x := range f.Field[y] {
			if f.Field[y][x].IsActive {
				r.renderFieldElement(
					screen,
					int32(fieldBiasX)+int32(x)*textureWidth,
					int32(fieldBiasY)+int32(y)*textureHeight,
					f.Field[y][x].Color)
			}
		}
	}
}

func (r *Render) renderFigure(screen *et.Image, f *field.Figure) {
	stateWithRightCoords := f.GetRightCoords()

	for i := range stateWithRightCoords.Coords {
		r.renderFieldElement(screen,
			int32(fieldBiasX)+int32(stateWithRightCoords.Coords[i].X)*textureWidth,
			int32(fieldBiasY)+int32(stateWithRightCoords.Coords[i].Y)*textureHeight,
			f.GetColor())
	}
}

func (r *Render) renderInfo(screen *et.Image) {
	if DrawFps {
		//rl.DrawFPS(2, 0) TODO
	}
}

func (r *Render) renderWindow(screen *et.Image) { // TODO

}

func addTextureToMap(key int, TexturePackName, name string) {
	var err error
	textures[key], _, err = drw.NewImageFromFile(fmt.Sprintf(pathToTexturesFormat, TexturePackName, name))
	if err != nil {
		log.Fatalln(err)
	}
}

func setConfigAndUploadTextures() {
	file, err := os.Open("client_side/render/render_config.json")
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
	TargetFps = config.TargetFps
	textureHeight = config.TextureHeight
	textureWidth = config.TextureWidth
	bonusPxWindowHeight = config.BonusPxWindowHeight
	bonusPxWindowWidth = config.BonusPxWindowWidth
	bonusPercentsWindowWidth = config.BonusPercentsWindowWidth
	bonusPercentsWindowHeight = config.BonusPercentsWindowHeight
	fieldBiasX = config.FieldBiasX
	fieldBiasY = config.FieldBiasY
	DrawFps = config.PrintFps

	windowX = int((int32(field.FieldX)*textureWidth)*(1+bonusPercentsWindowWidth/100)) + int(bonusPxWindowWidth+int32(fieldBiasX))
	windowY = int((int32(field.FieldY)*textureHeight)*(1+bonusPercentsWindowHeight/100)) + int(bonusPxWindowHeight+int32(fieldBiasY))

	textures = make(map[int]*et.Image)

	addTextureToMap(constants.ColorBlack, config.TexturePackName, "black")
	addTextureToMap(constants.ColorBlue, config.TexturePackName, "blue")
	addTextureToMap(constants.ColorGreen, config.TexturePackName, "green")
	addTextureToMap(constants.ColorOrange, config.TexturePackName, "orange")
	addTextureToMap(constants.ColorRed, config.TexturePackName, "red")
	addTextureToMap(constants.ColorYellow, config.TexturePackName, "yellow")
}

func init() {
	setConfigAndUploadTextures()
}
