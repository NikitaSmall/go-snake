package main

import (
	"math/rand"
	"time"
)

type gameState int

const slowestSpeed = 700 * time.Millisecond
const fastestSpeed = 60 * time.Millisecond

const (
	gameIntro gameState = iota
	gameStarted
	gamePaused
	gameOver
)

type bodyElement [2]int

type snake struct {
	body      []bodyElement
	direction int // from 0 to 3 up, right, down, left
}

type Game struct {
	board     [][]int
	state     gameState
	level     int
	snake     *snake
	turnTimer *time.Timer
}

func newSnake() *snake {
	body := make([]bodyElement, 3)
	for i := 0; i < 3; i++ {
		body[i] = [2]int{boardHeight / 2, boardWidth/2 + i}
	}

	return &snake{
		body:      body,
		direction: 0,
	}
}

func resetBoard() [][]int {
	board := make([][]int, boardHeight)
	for y := 0; y < boardHeight; y++ {
		board[y] = make([]int, boardWidth)

		for x := 0; x < boardWidth; x++ {
			board[y][x] = 0
		}
	}

	for y := 0; y < boardHeight; y++ {
		board[y][0] = 3
		board[y][boardWidth-1] = 3
	}

	for x := 0; x < boardWidth; x++ {
		board[0][x] = 3
		board[boardHeight-1][x] = 3
	}

	return board
}

func NewGame() *Game {
	g := new(Game)
	g.resetGame()
	return g
}

func (game *Game) resetGame() {
	game.snake = newSnake()
	game.board = resetBoard()

	game.level = 5
	game.state = gameIntro

	game.turnTimer = time.NewTimer(time.Duration(1000000 * time.Second))
	game.turnTimer.Stop()
}

func (g *Game) resetTurnTimer() {
	g.turnTimer.Reset(g.speed())
}

func (g *Game) speed() time.Duration {
	return slowestSpeed - fastestSpeed*time.Duration(g.level)
}

func (g *Game) play() {
	g.updateSnake()

	pos := g.snake.nextPos()
	if !g.checkPos(pos) {
		g.state = gameOver
	}

	g.resetTurnTimer()
}

func (s *snake) nextPos() [2]int {
	var pos [2]int

	switch s.direction {
	case 0:
		pos = [2]int{s.body[0][0] - 1, s.body[0][1]}
	case 1:
		pos = [2]int{s.body[0][0], s.body[0][1] + 1}
	case 2:
		pos = [2]int{s.body[0][0] + 1, s.body[0][1]}
	case 3:
		pos = [2]int{s.body[0][0], s.body[0][1] - 1}
	}

	return pos
}

func (g *Game) checkPos(pos [2]int) bool {
	cellValue := g.board[pos[0]][pos[1]]

	switch cellValue {
	case 0:
		g.snake.move(pos)
	case 1:
		g.board[pos[0]][pos[1]] = 2
		g.snake.grow(pos)
		g.placeFrog()
	default:
		return false
	}

	return true
}

func (s *snake) getBody() []bodyElement {
	return s.body
}

func (g *Game) updateSnake() {
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if g.board[y][x] == 2 {
				g.board[y][x] = 0
			}
		}
	}

	snakeBody := g.snake.getBody()
	for _, BodyEl := range snakeBody {
		g.board[BodyEl[0]][BodyEl[1]] = 2
	}
}

func (s *snake) move(pos bodyElement) {
	s.body = s.body[:len(s.body)-1]
	s.body = append([]bodyElement{pos}, s.body...)
}

func (s *snake) grow(pos bodyElement) {
	s.body = append([]bodyElement{pos}, s.body...)
}

func randPoint() (int, int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	y := r.Intn(boardHeight - 1)
	x := r.Intn(boardWidth - 1)

	return y, x
}

func (g *Game) placeFrog() {
	y, x := randPoint()

	for g.board[y][x] != 0 {
		y, x = randPoint()
	}

	g.board[y][x] = 1
}

func (s *snake) setDirection(dir int) {
	if s.direction%2 != dir%2 {
		s.direction = dir
	}
}

func (g *Game) moveLeft() {
	g.snake.setDirection(3)
}

func (g *Game) moveUp() {
	g.snake.setDirection(0)
}

func (g *Game) moveDown() {
	g.snake.setDirection(2)
}

func (g *Game) moveRight() {
	g.snake.setDirection(1)
}

func (g *Game) start() {
	switch g.state {
	case gameStarted:
		return
	case gamePaused:
		g.resume()
		return
	case gameOver:
		g.resetGame()
		fallthrough
	default:
		g.state = gameStarted
		g.placeFrog()
		g.resetTurnTimer()
	}
}

func (g *Game) pause() {
	switch g.state {
	case gameStarted:
		g.state = gamePaused
		g.turnTimer.Stop()
	case gamePaused:
		g.resume()
	}
}

func (g *Game) resume() {
	g.state = gameStarted
	g.play()
}
