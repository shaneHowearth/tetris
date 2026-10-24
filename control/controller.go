package control

import (
	"time"

	"github.com/shanehowearth/tetris/state"
)

type GameController struct {
	Board    *state.Board
	OnUpdate func()
	stopChan chan struct{}
}

func NewGameController(board *state.Board, onUpdate func()) *GameController {
	return &GameController{
		Board:    board,
		OnUpdate: onUpdate,
		stopChan: make(chan struct{}),
	}
}

func (c *GameController) StartGameLoop() {
	ticker := time.NewTicker(500 * time.Millisecond) // Gravity tick rate
	go func() {
		for {
			select {
			case <-ticker.C:
				if !c.Board.GameOver {
					c.Board.TickDown()
					c.OnUpdate()
				}
			case <-c.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *GameController) HandleInput(action string) {
	if c.Board.GameOver {
		return
	}

	switch action {
	case "LEFT":
		c.Board.MoveLeft()
	case "RIGHT":
		c.Board.MoveRight()
	case "ROTATE":
		c.Board.Rotate()
	case "DOWN":
		c.Board.TickDown()
	case "DROP":
		c.Board.HardDrop()
	}
	c.OnUpdate()
}

func (c *GameController) Restart() {
	*c.Board = *state.NewBoard(20, 10)
	c.OnUpdate()
}
