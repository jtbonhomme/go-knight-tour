package game

import (
	"image/color"
	"log"
	"time"

	"github.com/jtbonhomme/go-knight-tour/internal/knight"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BlinkFrameRate uint64 = 30
	Started               = iota
	Running
	GameWon
	GameLost
)

// Game manages all internal game mechanisms.
type Game struct {
	ScreenWidth       int
	ScreenHeight      int
	BackgroundColor   color.Color
	Knight            *knight.Knight
	state             int
	start             time.Time
	duration          time.Duration
	runResult         chan bool
	blinkFrameCounter uint64
	blink             bool
	stopChannel       chan struct{}
	slowMotion        int
	implementation    string
	slowMotionChange  chan int
	lastTime          time.Time
	debug             bool
}

// New creates a new game object.
func New(slowMotion int, implementation string, debug bool) *Game {
	g := &Game{
		ScreenWidth:      500,
		ScreenHeight:     500,
		BackgroundColor:  color.RGBA{0x0b, 0x0d, 0x00, 0xff},
		slowMotion:       slowMotion,
		implementation:   implementation,
		state:            Started,
		slowMotionChange: make(chan int),
		debug:            debug,
	}

	return g
}

// Run game loop.
func (g *Game) Run() error {

	ebiten.SetWindowSize(g.ScreenWidth, g.ScreenHeight)
	ebiten.SetWindowTitle("Knight Tour")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return ebiten.RunGame(g)
}

func (g *Game) Restart() {
	if g.state == Running {
		log.Println("restart game")
		close(g.stopChannel)
	}
	g.stopChannel = make(chan struct{})
	g.state = Running
	g.duration = 0
	g.start = time.Now()
	g.lastTime = time.Now()
	g.Knight = knight.New(g.slowMotion, g.implementation, g.slowMotionChange)
	g.runResult = g.Knight.Run(g.stopChannel)
	log.Printf("game: run")
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}
