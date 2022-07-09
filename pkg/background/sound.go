package background

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MusicController struct {
	rl.Music
	paused bool
}

const (
	bgMusicFile = "res/Avery's Farm Loopable.mp3"
)

func NewMusicController(paused bool) *MusicController {
	m := rl.LoadMusicStream(bgMusicFile)
	rl.PlayMusicStream(m)

	return &MusicController{
		Music:  m,
		paused: paused,
	}
}

func (mc *MusicController) Input() {
	if rl.IsKeyPressed(rl.KeyP) {
		mc.paused = !mc.paused
	}
}

func (mc *MusicController) Update(tick int) {
	rl.UpdateMusicStream(mc.Music)
	if mc.paused {
		rl.PauseMusicStream(mc.Music)
	} else {
		rl.ResumeMusicStream(mc.Music)
	}
}

func (mc *MusicController) Draw() {
}

func (mc *MusicController) Quit() {
	rl.UnloadMusicStream(mc.Music)
}
