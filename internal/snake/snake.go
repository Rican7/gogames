package snake

import (
	"image"
	"image/color"
	"strconv"

	"github.com/Rican7/gogames/gameengine/snake"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

// FrontEnd TODO
type FrontEnd struct {
	gameEngine *snake.GameEngine

	width  int
	height int

	statusAreaBounds             image.Rectangle
	playAreaBounds               image.Rectangle
	playAreaGameWidthMultiplier  uint
	playAreaGameHeightMultiplier uint
}

// NewFrontEnd TODO
func NewFrontEnd(gameEngine *snake.GameEngine, width int, height int) *FrontEnd {
	statusAreaBounds := image.Rect(0, 0, width, 20)
	playAreaBounds := image.Rect(0, statusAreaBounds.Max.Y, width, height)

	playAreaGameWidthMultiplier := uint(playAreaBounds.Dx()) / gameEngine.PlayAreaWidth()
	playAreaGameHeightMultiplier := uint(playAreaBounds.Dy()) / gameEngine.PlayAreaHeight()

	return &FrontEnd{
		gameEngine: gameEngine,

		width:  width,
		height: height,

		statusAreaBounds:             statusAreaBounds,
		playAreaBounds:               playAreaBounds,
		playAreaGameWidthMultiplier:  playAreaGameWidthMultiplier,
		playAreaGameHeightMultiplier: playAreaGameHeightMultiplier,
	}
}

// Run TODO
func (fe *FrontEnd) Run() error {
	ebiten.SetWindowSize(fe.width*2, fe.height*2)
	ebiten.SetWindowTitle("Snake")

	return ebiten.RunGame(fe)
}

// Layout TODO
func (fe *FrontEnd) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return fe.width, fe.height
}

// Update TODO
func (fe *FrontEnd) Update() error {
	gameStatus := fe.gameEngine.Status()

	if gameStatus == snake.StatusNew || gameStatus == snake.StatusPlaying {
		fe.handleInput()
		_ = fe.gameEngine.Tick()
	}

	return nil
}

// Draw TODO
func (fe *FrontEnd) Draw(screen *ebiten.Image) {
	screen.Fill(colorBlack)
	fe.drawStatusArea(screen)
	fe.drawPlayArea(screen)
}

func (fe *FrontEnd) handleInput() {
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

func (fe *FrontEnd) drawStatusArea(screen *ebiten.Image) {
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

	ebitenutil.DebugPrint(screen, strconv.Itoa(fe.gameEngine.Score()))
	// text.Draw(screen, "test", arcadeFont, screenWidth-len(scoreStr)*fontSize, fontSize, colorWhite)
}

func (fe *FrontEnd) drawPlayArea(screen *ebiten.Image) {
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
