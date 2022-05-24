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

	fieldBiasX uint32
	fieldBiasY uint32

	font fnt.Face
)

const pathToTexturesFormat = "client_side/render/textures/%s/texture_sq_%s.png"

func GetX() int {
	return windowX
}

func GetY() int {
	return windowY
}

func (r *Render) RenderAll(screen *et.Image, field *field.Field, figure *field.Figure, scores int, endGame bool) {
	r.renderFieldBackground(screen, r.FieldSize.X, r.FieldSize.Y)
	r.renderField(screen, field)
	if !endGame {
		r.renderFigure(screen, figure)
		r.renderInfo(screen, scores, field)
	} else {
		r.renderEndgameInfo(screen, scores, field)
	}
}

func (r *Render) renderEndgameInfo(screen *et.Image, scores int, f *field.Field) { //TODO
	text.Draw(screen, "Game over\nYour scores:"+strconv.Itoa(scores)+"\nPress ESC to exit", font, int(int32(fieldBiasX)+(textureWidth+1)*int32(f.GetSize().X)),
		int(fieldBiasY), color.RGBA{R: 255, G: 255, B: 255, A: 255})
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

func (r *Render) renderInfo(screen *et.Image, scores int, f *field.Field) {
	text.Draw(screen, strconv.Itoa(scores), font, int(int32(fieldBiasX)+(textureWidth+1)*int32(f.GetSize().X)),
		int(fieldBiasY), color.RGBA{R: 255, G: 255, B: 255, A: 255})
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
	textureHeight = config.TextureHeight
	textureWidth = config.TextureWidth
	bonusPxWindowHeight = config.BonusPxWindowHeight
	bonusPxWindowWidth = config.BonusPxWindowWidth
	bonusPercentsWindowWidth = config.BonusPercentsWindowWidth
	bonusPercentsWindowHeight = config.BonusPercentsWindowHeight
	fieldBiasX = config.FieldBiasX
	fieldBiasY = config.FieldBiasY

	windowX = int((int32(field.FieldX)*textureWidth)*(1+bonusPercentsWindowWidth/100)) + int(bonusPxWindowWidth+int32(fieldBiasX))
	windowY = int((int32(field.FieldY)*textureHeight)*(1+bonusPercentsWindowHeight/100)) + int(bonusPxWindowHeight+int32(fieldBiasY))

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
	setConfigAndUploadTextures()
}
