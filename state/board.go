package state

import (
	"math/rand"
	"time"
)

type BlockType int

const (
	BlockEmpty BlockType = iota
	BlockI               // Cyan
	BlockJ               // Blue
	BlockL               // Orange
	BlockO               // Yellow
	BlockS               // Green
	BlockT               // Purple
	BlockZ               // Red
)

type Point struct {
	Row int
	Col int
}

type Piece struct {
	Type  BlockType
	Shape [][]int // 1s and 0s defining current shape layout
	Row   int     // Top-left row position on grid
	Col   int     // Top-left col position on grid
}

type Board struct {
	Rows         int
	Cols         int
	Grid         [][]BlockType
	Current      *Piece
	Score        int
	LinesCleared int
	GameOver     bool
}

func NewBoard(rows, cols int) *Board {
	rand.Seed(time.Now().UnixNano())
	b := &Board{
		Rows: rows,
		Cols: cols,
		Grid: make([][]BlockType, rows),
	}
	for r := 0; r < rows; r++ {
		b.Grid[r] = make([]BlockType, cols)
	}
	b.SpawnPiece()
	return b
}

func (b *Board) SpawnPiece() {
	shapes := map[BlockType][][]int{
		BlockI: {{1, 1, 1, 1}},
		BlockJ: {{1, 0, 0}, {1, 1, 1}},
		BlockL: {{0, 0, 1}, {1, 1, 1}},
		BlockO: {{1, 1}, {1, 1}},
		BlockS: {{0, 1, 1}, {1, 1, 0}},
		BlockT: {{0, 1, 0}, {1, 1, 1}},
		BlockZ: {{1, 1, 0}, {0, 1, 1}},
	}

	t := BlockType(rand.Intn(7) + 1)
	shape := shapes[t]

	b.Current = &Piece{
		Type:  t,
		Shape: shape,
		Row:   0,
		Col:   (b.Cols - len(shape[0])) / 2,
	}

	if !b.isValidMove(b.Current.Shape, b.Current.Row, b.Current.Col) {
		b.GameOver = true
	}
}

func (b *Board) isValidMove(shape [][]int, newRow, newCol int) bool {
	for r := 0; r < len(shape); r++ {
		for c := 0; c < len(shape[r]); c++ {
			if shape[r][c] != 0 {
				gridR := newRow + r
				gridC := newCol + c

				if gridR < 0 || gridR >= b.Rows || gridC < 0 || gridC >= b.Cols {
					return false
				}
				if b.Grid[gridR][gridC] != BlockEmpty {
					return false
				}
			}
		}
	}
	return true
}

func (b *Board) MoveLeft() bool {
	if b.GameOver || b.Current == nil {
		return false
	}
	if b.isValidMove(b.Current.Shape, b.Current.Row, b.Current.Col-1) {
		b.Current.Col--
		return true
	}
	return false
}

func (b *Board) MoveRight() bool {
	if b.GameOver || b.Current == nil {
		return false
	}
	if b.isValidMove(b.Current.Shape, b.Current.Row, b.Current.Col+1) {
		b.Current.Col++
		return true
	}
	return false
}

func (b *Board) Rotate() bool {
	if b.GameOver || b.Current == nil {
		return false
	}
	rows := len(b.Current.Shape)
	cols := len(b.Current.Shape[0])

	// Rotate Matrix 90 Degrees Clockwise
	rotated := make([][]int, cols)
	for r := 0; r < cols; r++ {
		rotated[r] = make([]int, rows)
		for c := 0; c < rows; c++ {
			rotated[r][c] = b.Current.Shape[rows-1-c][r]
		}
	}

	if b.isValidMove(rotated, b.Current.Row, b.Current.Col) {
		b.Current.Shape = rotated
		return true
	}
	return false
}

func (b *Board) TickDown() bool {
	if b.GameOver || b.Current == nil {
		return false
	}

	if b.isValidMove(b.Current.Shape, b.Current.Row+1, b.Current.Col) {
		b.Current.Row++
		return true
	}

	// Lock piece in place
	b.lockPiece()
	b.clearLines()
	b.SpawnPiece()
	return false
}

func (b *Board) HardDrop() {
	if b.GameOver {
		return
	}
	for b.TickDown() {
	}
}

func (b *Board) lockPiece() {
	for r := 0; r < len(b.Current.Shape); r++ {
		for c := 0; c < len(b.Current.Shape[r]); c++ {
			if b.Current.Shape[r][c] != 0 {
				b.Grid[b.Current.Row+r][b.Current.Col+c] = b.Current.Type
			}
		}
	}
}

func (b *Board) clearLines() {
	cleared := 0
	for r := b.Rows - 1; r >= 0; r-- {
		full := true
		for c := 0; c < b.Cols; c++ {
			if b.Grid[r][c] == BlockEmpty {
				full = false
				break
			}
		}

		if full {
			cleared++
			for moveR := r; moveR > 0; moveR-- {
				copy(b.Grid[moveR], b.Grid[moveR-1])
			}
			b.Grid[0] = make([]BlockType, b.Cols)
			r++ // Check same row index again after shift
		}
	}

	b.LinesCleared += cleared
	b.Score += cleared * cleared * 100
}
