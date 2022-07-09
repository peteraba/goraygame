package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"gogame/pkg/background"
	"gogame/pkg/character"
)

type Actor interface {
	Input()
	Update(tick int)
	Draw()
	Quit()
}

type CamController interface {
	GetCamTarget() rl.Vector2
}

type Game struct {
	screenWidth   int32
	screenHeight  int32
	title         string
	fps           int32
	frameCount    int
	maxFrameCount int
	running       bool
	bgColor       rl.Color
	actors        []Actor
	player        *character.Player
	music         *background.MusicController
	flooring      *background.Flooring

	cam rl.Camera2D
}

var gameCount int

func New(screenWidth, screenHeight, fps int32, title string) *Game {
	if gameCount > 0 {
		panic("only one game can be handled at this point")
	}
	gameCount++

	return &Game{
		screenWidth:   screenWidth,
		screenHeight:  screenHeight,
		title:         title,
		fps:           fps,
		maxFrameCount: int(fps) * 60,
		running:       true,
		bgColor:       rl.NewColor(147, 211, 196, 255),
	}
}

func (g *Game) Input() {
	if rl.IsKeyDown(rl.KeyQ) {
		g.running = false
		return
	}

	for _, actor := range g.actors {
		actor.Input()
	}
}

func (g *Game) Update() {
	g.running = g.running && !rl.WindowShouldClose()

	for _, actor := range g.actors {
		actor.Update(g.frameCount)
	}

	g.cam.Target = g.player.GetCamTarget()

	g.frameCount++
	if g.frameCount > g.maxFrameCount {
		g.frameCount = 0
	}
}

func (g *Game) Render() {
	rl.BeginDrawing()
	defer rl.EndDrawing()
	rl.ClearBackground(g.bgColor)
	rl.BeginMode2D(g.cam)
	defer rl.EndMode2D()

	for _, actor := range g.actors {
		actor.Draw()
	}
}

func (g *Game) Init(playerSpeed float32, pausedMusic bool, zoom float32, flooringMapFile string) error {
	rl.InitWindow(g.screenWidth, g.screenHeight, g.title)
	rl.SetExitKey(0)
	rl.SetTargetFPS(g.fps)

	rl.InitAudioDevice()
	g.music = background.NewMusicController(pausedMusic)

	flooring, err := background.LoadFlooring(flooringMapFile)
	if err != nil {
		return fmt.Errorf("failed to initialize background. err: %w", err)
	}

	g.player = character.NewPlayer(playerSpeed, float32(g.screenWidth), float32(g.screenHeight))

	g.cam = rl.NewCamera2D(
		rl.NewVector2(float32(g.screenWidth)/2, float32(g.screenWidth)/2),
		g.player.GetCamTarget(),
		0.0,
		zoom,
	)

	g.flooring = &flooring
	g.actors = []Actor{g.music, g.flooring, g.player}

	return nil
}

func (g *Game) Quit() {
	for _, actor := range g.actors {
		actor.Quit()
	}
	rl.CloseWindow()
	rl.CloseAudioDevice()
}

func (g *Game) Run() {
	for g.running {
		g.Input()
		g.Update()
		g.Render()
	}
}
