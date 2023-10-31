package snake

import (
	"math/rand"
	"slices"
)

// Status TODO
type Status uint8

// TODO
const (
	StatusNew Status = iota
	StatusPlaying
	StatusLost
)

// Point TODO
type Point struct {
	X, Y uint
}

// Direction TODO
type Direction uint8

// TODO
const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

func (d Direction) isOpposite(dir Direction) bool {
	return (d == DirectionUp && dir == DirectionDown) ||
		(d == DirectionDown && dir == DirectionUp) ||
		(d == DirectionLeft && dir == DirectionRight) ||
		(d == DirectionRight && dir == DirectionLeft)
}

// GameEngine TODO
type GameEngine struct {
	status            Status
	playAreaWidth     uint
	playAreaHeight    uint
	score             int
	speed             uint
	tickThrottleCount uint

	foodLocation Point

	snakeBody               []Point // Head is first element, tail is last
	snakeDirection          Direction
	snakeDirectionLastMoved Direction
	snakeShouldElongate     bool
}

// NewGameEngine TODO
func NewGameEngine(playAreaWidth, playAreaHeight uint) *GameEngine {
	gameEngine := &GameEngine{
		playAreaWidth:  playAreaWidth,
		playAreaHeight: playAreaHeight,
	}

	gameEngine.NewGame()

	return gameEngine
}

// NewGame TODO
func (ge *GameEngine) NewGame() {
	ge.status = StatusNew
	ge.score = 0
	ge.speed = 1
	ge.tickThrottleCount = 0

	ge.foodLocation = Point{
		(ge.playAreaWidth / 4) * 3,
		ge.playAreaHeight / 2,
	}

	ge.snakeBody = []Point{
		Point{ge.playAreaWidth / 4, ge.playAreaHeight / 2},
		Point{(ge.playAreaWidth / 4) - 1, ge.playAreaHeight / 2},
	}
	ge.snakeDirection = DirectionRight
	ge.snakeDirectionLastMoved = DirectionRight
	ge.snakeShouldElongate = false
}

// Status TODO
func (ge *GameEngine) Status() Status {
	return ge.status
}

// PlayAreaWidth TODO
func (ge *GameEngine) PlayAreaWidth() uint {
	return ge.playAreaWidth
}

// PlayAreaHeight TODO
func (ge *GameEngine) PlayAreaHeight() uint {
	return ge.playAreaHeight
}

// Score TODO
func (ge *GameEngine) Score() int {
	return ge.score
}

// Speed TODO
func (ge *GameEngine) Speed() uint {
	return ge.speed
}

// FoodLocation TODO
func (ge *GameEngine) FoodLocation() Point {
	return ge.foodLocation
}

// SnakeBody TODO
func (ge *GameEngine) SnakeBody() []Point {
	return ge.snakeBody
}

// UpdateDirection TODO
func (ge *GameEngine) UpdateDirection(dir Direction) {
	// Don't update the direction if it's the same or opposite of it's current
	if ge.snakeDirection == dir || ge.snakeDirection.isOpposite(dir) ||
		ge.snakeDirectionLastMoved == dir || ge.snakeDirectionLastMoved.isOpposite(dir) {
		return
	}

	ge.snakeDirection = dir
}

// Tick TODO
func (ge *GameEngine) Tick() Status {
	if ge.status == StatusLost || !ge.tickThrottle() {
		return ge.status
	}

	if ge.status == StatusNew {
		ge.status = StatusPlaying
	}

	validMove := ge.moveSnake()

	if !validMove {
		ge.status = StatusLost
		return ge.status
	}

	if ge.isInScorePosition() {
		ge.score++

		// Increase the speed every 5 points
		if ge.score%5 == 0 {
			ge.speed++
		}

		ge.placeFood()
		ge.snakeShouldElongate = true
	}

	return ge.status
}

func (ge *GameEngine) tickThrottle() bool {
	if ge.tickThrottleCount == 0 {
		ge.tickThrottleCount = 10 - (ge.speed - 1)

		return true
	}

	ge.tickThrottleCount--

	return false
}

func (ge *GameEngine) moveSnake() bool {
	next := ge.snakeBody[0]
	direction := ge.snakeDirection

	switch direction {
	case DirectionUp:
		next = Point{next.X, next.Y - 1}
	case DirectionDown:
		next = Point{next.X, next.Y + 1}
	case DirectionLeft:
		next = Point{next.X - 1, next.Y}
	case DirectionRight:
		next = Point{next.X + 1, next.Y}
	}

	if ge.wouldMoveCauseLoss(next) {
		return false
	}

	ge.snakeBody = slices.Insert(ge.snakeBody, 0, next)
	ge.snakeDirectionLastMoved = direction

	if !ge.snakeShouldElongate {
		// Remove tail
		ge.snakeBody = ge.snakeBody[:len(ge.snakeBody)-1]
	}
	ge.snakeShouldElongate = false

	return true
}

func (ge *GameEngine) wouldMoveCauseLoss(next Point) bool {
	if ge.status == StatusLost {
		return true
	}

	// Check for out-of-bounds
	if next.X < 0 || next.Y < 0 ||
		next.X > ge.playAreaWidth-1 || next.Y > ge.playAreaHeight-1 {
		return true
	}

	// Check for intersection of next and other body parts
	for _, point := range ge.snakeBody {
		if next.X == point.X && next.Y == point.Y {
			return true
		}
	}

	return false
}

func (ge *GameEngine) isInScorePosition() bool {
	if ge.status == StatusLost {
		return false
	}

	head := ge.snakeBody[0]

	return head.X == ge.foodLocation.X && head.Y == ge.foodLocation.Y
}

func (ge *GameEngine) placeFood() {
	var newLocation Point
	validLocation := false

	for !validLocation {
		newLocation = Point{
			uint(rand.Intn(int(ge.playAreaWidth))),
			uint(rand.Intn(int(ge.playAreaHeight))),
		}

		validLocation = true

		// Check that it's not the existing location
		if newLocation.X == ge.foodLocation.X || newLocation.Y == ge.foodLocation.Y {
			validLocation = false
			continue
		}

		// Check that it's not ON the snake body
		if slices.Contains(ge.snakeBody, newLocation) {
			validLocation = false
			continue
		}
	}

	ge.foodLocation = newLocation
}
