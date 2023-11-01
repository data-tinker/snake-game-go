package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	width  = 20
	height = 10
)

type point struct{ x, y int }

var directions = map[string]point{
	"up":    {0, -1},
	"down":  {0, 1},
	"left":  {-1, 0},
	"right": {1, 0},
}

type snakeGame struct {
	snake         []point
	food          point
	dir           point
	nextDir       point
	score         int
	snakeBodyChar rune
	foodChar      rune
	spaceChar     rune
}

func newGame() *snakeGame {
	snake := []point{{width / 2, height / 2}}
	food := point{width / 4, height / 4}
	return &snakeGame{
		snake:         snake,
		food:          food,
		dir:           directions["right"],
		nextDir:       directions["right"],
		score:         0,
		snakeBodyChar: '■',
		foodChar:      '●',
		spaceChar:     ' ',
	}
}

func (g *snakeGame) moveSnake() {
	g.dir = g.nextDir
	head := g.snake[0]
	newHead := point{head.x + g.nextDir.x, head.y + g.nextDir.y}

	// Check if snake has hit the wall
	if newHead.x < 0 || newHead.x >= width || newHead.y < 0 || newHead.y >= height {
		g.endGame()
		return
	}

	// Check if snake has hit itself
	for _, segment := range g.snake {
		if newHead == segment {
			g.endGame()
			return
		}
	}

	// Move snake
	g.snake = append([]point{newHead}, g.snake...)
	if newHead == g.food {
		g.score++
		g.placeFood()
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}
}

func (g *snakeGame) placeFood() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		foodCandidate := point{r.Intn(width), r.Intn(height)}
		occupied := false
		for _, segment := range g.snake {
			if foodCandidate == segment {
				occupied = true
				break
			}
		}
		if !occupied {
			g.food = foodCandidate
			break
		}
	}
}

func (g *snakeGame) endGame() {
	fmt.Printf("\nGame Over! Your score is: %d\nPress Enter to exit...", g.score)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(0)
}

func (g *snakeGame) changeDirection() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			var newDir point
			switch ev.Key {
			case termbox.KeyArrowUp:
				newDir = directions["up"]
			case termbox.KeyArrowDown:
				newDir = directions["down"]
			case termbox.KeyArrowLeft:
				newDir = directions["left"]
			case termbox.KeyArrowRight:
				newDir = directions["right"]
			case termbox.KeyCtrlC:
				termbox.Close()
				os.Exit(0)
			}
			// Check if the new direction is not the opposite of the current direction
			if newDir.x != -g.dir.x || newDir.y != -g.dir.y {
				g.nextDir = newDir
			}
		}
	}
}

func (g *snakeGame) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			point := point{x, y}
			cell := termbox.Cell{Fg: termbox.ColorDefault, Bg: termbox.ColorDefault}

			switch {
			case point == g.food:
				cell.Ch = g.foodChar
			case point == g.snake[0]:
				cell.Ch = g.snakeBodyChar
			default:
				isBody := false
				for _, segment := range g.snake[1:] {
					if point == segment {
						cell.Ch = g.snakeBodyChar
						isBody = true
						break
					}
				}
				if !isBody {
					cell.Ch = g.spaceChar
				}
			}

			termbox.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
		}
	}

	scoreText := fmt.Sprintf("Score: %d", g.score)
	for i, ch := range scoreText {
		termbox.SetCell(i, height, ch, termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.Flush()
}

func clearScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
}

func randInt(max int) int {
	return int(time.Now().UnixNano() % int64(max))
}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	game := newGame()
	go game.changeDirection()
	ticker := time.NewTicker(time.Second / 2)
	for {
		select {
		case <-ticker.C:
			game.moveSnake()
			game.draw()
		}
	}
}
