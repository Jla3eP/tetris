package mainGame

import (
	"github.com/Jla3eP/tetris/client_side/field"
	"github.com/Jla3eP/tetris/client_side/render"
	et "github.com/hajimehoshi/ebiten/v2"
)

type MainGame struct {
	*render.Render
	*field.Field
}

func NewGame() *MainGame {
	mg := &MainGame{}
	mg.Field = field.NewField()
	mg.Render = createRenderObj(mg.GetSize())
	return mg
}

func (mg *MainGame) Update() error {
	return nil
}

func (mg *MainGame) Draw(screen *et.Image) {
	mg.RenderAll(screen)
}

func (mg *MainGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return render.GetX(), render.GetY()
}

//
//func (mg *MainGame) StartGame() {
//	rd := createRenderData()
//
//	rdCh := make(chan *render.RenderData, 1)
//	endCh := make(chan struct{})
//	eventCh := make(chan constants.Event, 0)
//	eventCtx, cancel := context.WithCancel(context.Background())
//	eventCtx = eventCtx
//	endCh = endCh
//	defer cancel()
//
//	go func(eventCh <-chan constants.Event) {
//
//		ev := constants.Event(0)
//		for {
//			select {
//			case ev = <-eventCh:
//				fmt.Printf("%v", ev)
//			}
//
//		}
//	}(eventCh)
//
//	ticker := time.NewTicker(250 * time.Millisecond)
//	go func() {
//		for {
//			rdCh <- &rd
//
//			<-ticker.C
//			_ = rd.Field.TryRotateFigure(rd.Figure)
//		}
//	}()
//}
//
//func createRenderData() render.RenderData {
//	figures := field.GetFigures()
//	rd := render.RenderData{
//		Field:  &field.Field{Field: make([][]field.FieldPart, 20)},
//		Figure: &figures[rand.Int()%len(figures)],
//	}
//	for i := 0; i < len(rd.Field.Field); i++ {
//		rd.Field.Field[i] = make([]field.FieldPart, 10)
//	}
//	return rd
//}
//
func createRenderObj(fieldSize field.Coords2) *render.Render {
	return &render.Render{FieldSize: fieldSize}
}
