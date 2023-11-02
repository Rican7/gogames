package snake

import (
	"image"
	"image/color"
	"strconv"

	"golang.org/x/image/font"

	"github.com/Rican7/gogames/gameengine/snake"
	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	// Colors
	colorBlack       = color.Black
	colorWhite       = color.White
	colorGrayLighter = color.Gray{0xe0}
	colorGrayLight   = color.Gray{0xc0}
	colorGray        = color.Gray{0x80}
	colorGrayDark    = color.Gray{0x40}
	colorGrayDarker  = color.Gray{0x20}
	colorBlue        = color.RGBA{0x0, 0x0, 0xff, 0xff}
)

// GUIFrontEnd is a Graphical User Interface Front End to the Snake game.
type GUIFrontEnd struct {
	gameEngine *snake.GameEngine

	width  int
	height int

	statusAreaBounds             image.Rectangle
	playAreaBounds               image.Rectangle
	playAreaGameWidthMultiplier  uint
	playAreaGameHeightMultiplier uint
}

// NewFrontEnd creates a new GUIFrontEnd with the given game engine and
// un-scaled GUI dimensions.
func NewFrontEnd(gameEngine *snake.GameEngine, width int, height int) *GUIFrontEnd {
	statusAreaBounds := image.Rect(0, 0, width, 20)
	playAreaBounds := image.Rect(0, statusAreaBounds.Max.Y, width, height)

	playAreaGameWidthMultiplier := uint(playAreaBounds.Dx()) / gameEngine.PlayAreaWidth()
	playAreaGameHeightMultiplier := uint(playAreaBounds.Dy()) / gameEngine.PlayAreaHeight()

	return &GUIFrontEnd{
		gameEngine: gameEngine,

		width:  width,
		height: height,

		statusAreaBounds:             statusAreaBounds,
		playAreaBounds:               playAreaBounds,
		playAreaGameWidthMultiplier:  playAreaGameWidthMultiplier,
		playAreaGameHeightMultiplier: playAreaGameHeightMultiplier,
	}
}

// Run runs the game's GUI Front End. It returns an error if any occurred.
func (fe *GUIFrontEnd) Run() error {
	ebiten.SetWindowSize(fe.width*2, fe.height*2)
	ebiten.SetWindowTitle("Snake")

	return ebiten.RunGame(fe)
}

// Layout takes the GUI window dimensions and returns the dimensions of the
// GUI's logical layout.
func (fe *GUIFrontEnd) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return fe.width, fe.height
}

// Update performs logical updates for the game to function, returning any
// errors caused.
func (fe *GUIFrontEnd) Update() error {
	gameStatus := fe.gameEngine.Status()

	switch gameStatus {
	case snake.StatusNew, snake.StatusLost:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			fe.gameEngine.PlayNew()
		}
	case snake.StatusPlaying:
		fe.handlePlayInput()
		fe.gameEngine.Tick()
	}

	return nil
}

// Draw takes a screen and draws the GUI elements to the screen.
func (fe *GUIFrontEnd) Draw(screen *ebiten.Image) {
	screen.Fill(colorBlack)
	fe.drawStatusArea(screen)
	fe.drawPlayArea(screen)
}

// handlePlayInput handles the input used during gameplay.
func (fe *GUIFrontEnd) handlePlayInput() {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowUp):
		fe.gameEngine.UpdateDirection(snake.DirectionUp)
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowDown):
		fe.gameEngine.UpdateDirection(snake.DirectionDown)
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft):
		fe.gameEngine.UpdateDirection(snake.DirectionLeft)
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowRight):
		fe.gameEngine.UpdateDirection(snake.DirectionRight)
	}
}

// drawStatusArea takes a screen and draws the status area to it.
func (fe *GUIFrontEnd) drawStatusArea(screen *ebiten.Image) {
	statusAreaBoundsCenterX := fe.statusAreaBounds.Min.X + (fe.statusAreaBounds.Dx() / 2)
	statusAreaBoundsCenterY := fe.statusAreaBounds.Min.Y + (fe.statusAreaBounds.Dy() / 2)

	vector.StrokeLine(
		screen,
		float32(fe.statusAreaBounds.Min.X),
		float32(fe.statusAreaBounds.Max.Y),
		float32(fe.statusAreaBounds.Max.X),
		float32(fe.statusAreaBounds.Max.Y),
		1,
		colorGrayDark,
		false,
	)

	scoreStr := strconv.Itoa(fe.gameEngine.Score())
	scoreBounds, _ := font.BoundString(bitmapfont.Face, scoreStr)

	text.Draw(
		screen,
		scoreStr,
		bitmapfont.Face,
		statusAreaBoundsCenterX-(scoreBounds.Max.X.Floor()/2),
		statusAreaBoundsCenterY+scoreBounds.Max.Y.Floor(),
		colorWhite,
	)
}

// drawPlayArea takes a screen and draws the play area to it.
func (fe *GUIFrontEnd) drawPlayArea(screen *ebiten.Image) {
	switch fe.gameEngine.Status() {
	case snake.StatusNew:
		fe.drawNewGame(screen)
	case snake.StatusPlaying:
		fe.drawRunningGame(screen)
	case snake.StatusLost:
		fe.drawGameOver(screen)
	}
}

// drawNewGame takes a screen and draws the new game screen graphics to it.
func (fe *GUIFrontEnd) drawNewGame(screen *ebiten.Image) {
	playAreaBoundsCenterX := fe.playAreaBounds.Min.X + (fe.playAreaBounds.Dx() / 2)
	playAreaBoundsCenterY := fe.playAreaBounds.Min.Y + (fe.playAreaBounds.Dy() / 2)
	fontSize := bitmapfont.Face.Metrics().Height.Ceil()
	newGameStr := "Snake"
	playAgainStr := "Press Enter to play"
	newGameBounds, _ := font.BoundString(bitmapfont.Face, newGameStr)
	playAgainBounds, _ := font.BoundString(bitmapfont.Face, playAgainStr)
	messageHeight := fontSize * 3 // 3 Lines: 2 messages with a blank between

	vector.DrawFilledRect(
		screen,
		float32(fe.playAreaBounds.Min.X),
		float32(fe.playAreaBounds.Min.Y),
		float32(fe.playAreaBounds.Dx()),
		float32(fe.playAreaBounds.Dy()),
		colorBlack,
		false,
	)

	text.Draw(
		screen,
		newGameStr,
		bitmapfont.Face,
		playAreaBoundsCenterX-(newGameBounds.Max.X.Floor()/2),
		playAreaBoundsCenterY- // Area center
			(messageHeight/2)+ // Entire combined message height center,
			newGameBounds.Max.Y.Floor(),
		colorWhite,
	)

	text.Draw(
		screen,
		playAgainStr,
		bitmapfont.Face,
		playAreaBoundsCenterX-(playAgainBounds.Max.X.Floor()/2),
		playAreaBoundsCenterY- // Area center
			(messageHeight/2)+ // Entire combined message height center
			(fontSize*2)+ // Font height offset based on previous lines
			playAgainBounds.Max.Y.Floor(),
		colorWhite,
	)
}

// drawRunningGame takes a screen and draws the running game graphics to it.
func (fe *GUIFrontEnd) drawRunningGame(screen *ebiten.Image) {
	foodLocation := fe.gameEngine.FoodLocation()
	snakeBody := fe.gameEngine.SnakeBody()

	vector.DrawFilledRect(
		screen,
		float32(fe.playAreaBounds.Min.X),
		float32(fe.playAreaBounds.Min.Y),
		float32(fe.playAreaBounds.Dx()),
		float32(fe.playAreaBounds.Dy()),
		colorGrayDarker,
		false,
	)

	vector.DrawFilledRect(
		screen,
		float32(foodLocation.X*fe.playAreaGameWidthMultiplier)+float32(fe.playAreaBounds.Min.X),
		float32(foodLocation.Y*fe.playAreaGameHeightMultiplier)+float32(fe.playAreaBounds.Min.Y),
		float32(fe.playAreaGameWidthMultiplier),
		float32(fe.playAreaGameHeightMultiplier),
		colorWhite,
		false,
	)

	for _, point := range snakeBody {
		vector.DrawFilledRect(
			screen,
			float32(point.X*fe.playAreaGameWidthMultiplier)+float32(fe.playAreaBounds.Min.X),
			float32(point.Y*fe.playAreaGameHeightMultiplier)+float32(fe.playAreaBounds.Min.Y),
			float32(fe.playAreaGameWidthMultiplier),
			float32(fe.playAreaGameHeightMultiplier),
			colorBlue,
			false,
		)
	}
}

// drawGameOver takes a screen and draws the game over screen graphics to it.
func (fe *GUIFrontEnd) drawGameOver(screen *ebiten.Image) {
	playAreaBoundsCenterX := fe.playAreaBounds.Min.X + (fe.playAreaBounds.Dx() / 2)
	playAreaBoundsCenterY := fe.playAreaBounds.Min.Y + (fe.playAreaBounds.Dy() / 2)
	fontSize := bitmapfont.Face.Metrics().Height.Ceil()
	gameOverStr := "Game Over"
	playAgainStr := "Press Enter to play again"
	gameOverBounds, _ := font.BoundString(bitmapfont.Face, gameOverStr)
	playAgainBounds, _ := font.BoundString(bitmapfont.Face, playAgainStr)
	messageHeight := fontSize * 3 // 3 Lines: 2 messages with a blank between

	vector.DrawFilledRect(
		screen,
		float32(fe.playAreaBounds.Min.X),
		float32(fe.playAreaBounds.Min.Y),
		float32(fe.playAreaBounds.Dx()),
		float32(fe.playAreaBounds.Dy()),
		colorBlack,
		false,
	)

	text.Draw(
		screen,
		gameOverStr,
		bitmapfont.Face,
		playAreaBoundsCenterX-(gameOverBounds.Max.X.Floor()/2),
		playAreaBoundsCenterY- // Area center
			(messageHeight/2)+ // Entire combined message height center,
			gameOverBounds.Max.Y.Floor(),
		colorWhite,
	)

	text.Draw(
		screen,
		playAgainStr,
		bitmapfont.Face,
		playAreaBoundsCenterX-(playAgainBounds.Max.X.Floor()/2),
		playAreaBoundsCenterY- // Area center
			(messageHeight/2)+ // Entire combined message height center
			(fontSize*2)+ // Font height offset based on previous lines
			playAgainBounds.Max.Y.Floor(),
		colorWhite,
	)
}
