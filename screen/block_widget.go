package screen

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/tetris/state"
)

type BlockWidget struct {
	widget.BaseWidget
	BlockType state.BlockType
}

func NewBlockWidget(bType state.BlockType) *BlockWidget {
	bw := &BlockWidget{BlockType: bType}
	bw.ExtendBaseWidget(bw)
	return bw
}

func (bw *BlockWidget) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(getBlockColor(bw.BlockType))
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 1
	border.StrokeColor = color.RGBA{R: 30, G: 30, B: 30, A: 255}

	return &blockWidgetRenderer{
		bw:      bw,
		rect:    rect,
		border:  border,
		objects: []fyne.CanvasObject{rect, border},
	}
}

func getBlockColor(b state.BlockType) color.Color {
	switch b {
	case state.BlockI:
		return color.RGBA{R: 0, G: 240, B: 240, A: 255} // Cyan
	case state.BlockJ:
		return color.RGBA{R: 0, G: 30, B: 240, A: 255} // Blue
	case state.BlockL:
		return color.RGBA{R: 240, G: 160, B: 0, A: 255} // Orange
	case state.BlockO:
		return color.RGBA{R: 240, G: 240, B: 0, A: 255} // Yellow
	case state.BlockS:
		return color.RGBA{R: 0, G: 240, B: 0, A: 255} // Green
	case state.BlockT:
		return color.RGBA{R: 160, G: 0, B: 240, A: 255} // Purple
	case state.BlockZ:
		return color.RGBA{R: 240, G: 0, B: 0, A: 255} // Red
	default:
		return color.RGBA{R: 15, G: 15, B: 20, A: 255} // Grid background
	}
}

type blockWidgetRenderer struct {
	bw      *BlockWidget
	rect    *canvas.Rectangle
	border  *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *blockWidgetRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
	r.border.Resize(size)
}

func (r *blockWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(24, 24)
}

func (r *blockWidgetRenderer) Refresh() {
	r.rect.FillColor = getBlockColor(r.bw.BlockType)
	r.rect.Refresh()
}

func (r *blockWidgetRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *blockWidgetRenderer) Destroy()                     {}
