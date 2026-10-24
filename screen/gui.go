package screen

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/tetris/control"
	"github.com/shanehowearth/tetris/state"
)

type TetrisView struct {
	App        fyne.App
	Win        fyne.Window
	Board      *state.Board
	Ctrl       *control.GameController
	GridBox    *fyne.Container
	ScoreLabel *widget.Label
	Status     *widget.Label
}

func StartApp() {
	myApp := app.New()
	win := myApp.NewWindow("Go Tetris")

	board := state.NewBoard(20, 10) // Standard 20x10 Tetris Grid
	view := &TetrisView{
		App:   myApp,
		Win:   win,
		Board: board,
	}

	ctrl := control.NewGameController(board, view.RefreshUI)
	view.Ctrl = ctrl

	view.setupUI()
	ctrl.StartGameLoop()
	win.ShowAndRun()
}

func (v *TetrisView) setupUI() {
	v.ScoreLabel = widget.NewLabel("Score: 0")
	v.Status = widget.NewLabel("Playing")

	// Grid layout
	v.GridBox = container.NewGridWithColumns(v.Board.Cols)
	v.rebuildGrid()

	sidebar := container.NewVBox(
		widget.NewLabelWithStyle("TETRIS", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		v.ScoreLabel,
		v.Status,
		widget.NewSeparator(),
		widget.NewButton("Restart", func() {
			v.Ctrl.Restart()
		}),
	)

	content := container.NewHBox(v.GridBox, widget.NewSeparator(), sidebar)
	v.Win.SetContent(container.NewPadded(content))

	// Listen to Canvas Key Presses
	v.Win.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyLeft:
			v.Ctrl.HandleInput("LEFT")
		case fyne.KeyRight:
			v.Ctrl.HandleInput("RIGHT")
		case fyne.KeyUp:
			v.Ctrl.HandleInput("ROTATE")
		case fyne.KeyDown:
			v.Ctrl.HandleInput("DOWN")
		case fyne.KeySpace:
			v.Ctrl.HandleInput("DROP")
		}
	})
}

func (v *TetrisView) rebuildGrid() {
	v.GridBox.Objects = nil

	// Composite temporary grid combining locked grid + active piece
	renderGrid := make([][]state.BlockType, v.Board.Rows)
	for r := 0; r < v.Board.Rows; r++ {
		renderGrid[r] = make([]state.BlockType, v.Board.Cols)
		copy(renderGrid[r], v.Board.Grid[r])
	}

	if v.Board.Current != nil && !v.Board.GameOver {
		p := v.Board.Current
		for r := 0; r < len(p.Shape); r++ {
			for c := 0; c < len(p.Shape[r]); c++ {
				if p.Shape[r][c] != 0 {
					gridR := p.Row + r
					gridC := p.Col + c
					if gridR >= 0 && gridR < v.Board.Rows && gridC >= 0 && gridC < v.Board.Cols {
						renderGrid[gridR][gridC] = p.Type
					}
				}
			}
		}
	}

	for r := 0; r < v.Board.Rows; r++ {
		for c := 0; c < v.Board.Cols; c++ {
			v.GridBox.Add(NewBlockWidget(renderGrid[r][c]))
		}
	}
}

func (v *TetrisView) RefreshUI() {
	v.ScoreLabel.SetText(fmt.Sprintf("Score: %d", v.Board.Score))
	if v.Board.GameOver {
		v.Status.SetText("GAME OVER")
	} else {
		v.Status.SetText("Playing")
	}
	v.rebuildGrid()
	v.GridBox.Refresh()
}
