package character

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	speed      float32
	moving     bool
	dir        int
	up         bool
	down       bool
	right      bool
	left       bool
	frame      int
	texture    *rl.Texture2D
	src        rl.Rectangle
	dest       rl.Rectangle
	destVector rl.Vector2
}

const (
	playerTextureFile   = "res/Sprout Lands - Sprites - Basic pack/Characters/Basic Charakter Spritesheet.png"
	playerTextureWidth  = 48
	playerTextureHeight = 48
	wtf                 = 60
)

func NewPlayer(speed float32, screenWidth, screenHeight float32) *Player {
	texture := rl.LoadTexture(playerTextureFile)

	var (
		dw float32 = 200
		dh float32 = 200
	)

	log.Println(screenWidth, screenHeight, dw, dh)

	return &Player{
		speed:      speed,
		texture:    &texture,
		src:        rl.NewRectangle(0, 0, playerTextureWidth, playerTextureHeight),
		dest:       rl.NewRectangle(dw, dh, wtf, wtf),
		destVector: rl.NewVector2(wtf, wtf),
	}
}

func (p *Player) Input() {
	p.moving = false
	p.up, p.down, p.right, p.left = false, false, false, false
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		p.dest.Y -= p.speed
		p.moving = true
		p.up = true
		p.dir = 1
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		p.dest.Y += p.speed
		p.moving = true
		p.down = true
		p.dir = 0
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		p.dest.X -= p.speed
		p.moving = true
		p.left = true
		p.dir = 2
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		p.dest.X += p.speed
		p.moving = true
		p.right = true
		p.dir = 3
	}
}

func (p *Player) Update(tick int) {
	if p.moving {
		p.updateMoving(tick)
	} else if tick%45 == 0 {
		p.updateStationary(tick)
	}

	p.src.X = p.src.Width * float32(p.frame)
	p.src.Y = p.src.Height * float32(p.dir)

	//if tick%10 == 0 {
	//	log.Printf("player. src: %#v, dest: %#v, vector: %#v, tick: %d, frame: %d, moving: %t", p.src, p.dest, p.destVector, tick, p.frame, p.moving)
	//}
}

func (p *Player) updateMoving(tick int) {
	if p.up {
		p.dest.Y -= p.speed
	}
	if p.down {
		p.dest.Y += p.speed
	}
	if p.right {
		p.dest.X += p.speed
	}
	if p.left {
		p.dest.X -= p.speed
	}

	// Animate every 8th frame while moving (~.13s)
	if tick%8 == 0 {
		p.frame++
	}

	// During moving we use all 4 frames
	if p.frame > 3 {
		p.frame = 0
	}
}

func (p *Player) updateStationary(tick int) {
	// Animate every 45th frame while stationary (~.75s)
	if tick%45 == 0 {
		p.frame++
	}

	// During moving we use the first 2 frames
	if p.frame > 1 {
		p.frame = 0
	}
}

func (p *Player) Draw() {
	// TODO: Fix positions, this is kind of a hack for cam target
	rl.DrawTexturePro(*p.texture, p.src, rl.NewRectangle(p.dest.X, p.dest.Y-120, p.dest.Width, p.dest.Height), p.destVector, 0, rl.White)
}

func (p *Player) Quit() {
	rl.UnloadTexture(*p.texture)
}

func (p *Player) GetCamTarget() rl.Vector2 {
	return rl.NewVector2(p.dest.X-(p.dest.Width/2), p.dest.Y-(p.dest.Height/2))
}
