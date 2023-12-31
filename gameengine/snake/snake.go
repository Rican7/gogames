// Package snake provides a logical game-engine for a two-dimensional game in
// the classic "snake" genre.
//
// See https://en.wikipedia.org/wiki/Snake_(video_game_genre)
package snake

import (
	"math/rand"
	"slices"
)

// Status defines the running state status of the game.
type Status uint8

// The valid running state Status values.
const (
	StatusNew Status = iota
	StatusPlaying
	StatusLost
)

// Point is a coordinate pair, representing a point on a 2D plane.
type Point struct {
	X, Y uint
}

// Direction defines a 2D facing or moving direction.
type Direction uint8

// The valid Direction values.
const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

// isOpposite takes another direction and returns true if it's the opposite.
func (d Direction) isOpposite(dir Direction) bool {
	return (d == DirectionUp && dir == DirectionDown) ||
		(d == DirectionDown && dir == DirectionUp) ||
		(d == DirectionLeft && dir == DirectionRight) ||
		(d == DirectionRight && dir == DirectionLeft)
}

// GameEngine is the game engine for the game Snake.
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

// NewGameEngine creates a new GameEngine with the given play area dimensions,
// defining the 2D plane on which to play.
func NewGameEngine(playAreaWidth, playAreaHeight uint) *GameEngine {
	gameEngine := &GameEngine{
		playAreaWidth:  playAreaWidth,
		playAreaHeight: playAreaHeight,
	}

	gameEngine.NewGame()

	return gameEngine
}

// NewGame (re)initializes the game engine's state to that of a new, fresh game.
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

// Status returns the status of the game.
func (ge *GameEngine) Status() Status {
	return ge.status
}

// PlayAreaWidth returns the width of the play area.
func (ge *GameEngine) PlayAreaWidth() uint {
	return ge.playAreaWidth
}

// PlayAreaHeight returns the height of the play area.
func (ge *GameEngine) PlayAreaHeight() uint {
	return ge.playAreaHeight
}

// Score returns the current score.
func (ge *GameEngine) Score() int {
	return ge.score
}

// Speed returns the current speed of the game.
func (ge *GameEngine) Speed() uint {
	return ge.speed
}

// FoodLocation returns the location of the current "food" piece.
func (ge *GameEngine) FoodLocation() Point {
	return ge.foodLocation
}

// SnakeBody returns the location of the snake's body parts, as a slice of 2D
// points, with the head being the first element and the tail being the last.
func (ge *GameEngine) SnakeBody() []Point {
	return ge.snakeBody
}

// UpdateDirection updates the direction of the snake.
func (ge *GameEngine) UpdateDirection(dir Direction) {
	// Don't update the direction if it's the same or opposite of it's current
	if ge.snakeDirection == dir || ge.snakeDirection.isOpposite(dir) ||
		ge.snakeDirectionLastMoved == dir || ge.snakeDirectionLastMoved.isOpposite(dir) {
		return
	}

	ge.snakeDirection = dir
}

// PlayNew initializes the game state to a new game and sets the running state
// to StatusPlaying.
func (ge *GameEngine) PlayNew() {
	ge.NewGame()
	ge.status = StatusPlaying
}

// Tick performs a single iteration of the main game loop logic and returns the
// status of the game.
func (ge *GameEngine) Tick() Status {
	if ge.status != StatusPlaying || !ge.tickThrottle() {
		return ge.status
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

// tickThrottle is a rudimentary throttling mechanism for the game loop. It
// returns true if the game loop logic should execute and false otherwise.
func (ge *GameEngine) tickThrottle() bool {
	if ge.tickThrottleCount == 0 {
		ge.tickThrottleCount = 10 - (ge.speed - 1)

		return true
	}

	ge.tickThrottleCount--

	return false
}

// moveSnake moves the snake in the current direction and returns whether the
// move was valid (an invalid move should cause a game loss).
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

// wouldMoveCauseLoss returns whether moving the snake to the given point should
// cause a game loss or not.
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

// isInScorePosition returns whether the snake is currently in a scoring
// position or not.
func (ge *GameEngine) isInScorePosition() bool {
	if ge.status == StatusLost {
		return false
	}

	head := ge.snakeBody[0]

	return head.X == ge.foodLocation.X && head.Y == ge.foodLocation.Y
}

// placeFood places the "food" piece onto a new location on the board, randomly,
// making sure that the placement is valid.
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
