package background

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	grass  = "g"
	hills  = "l"
	fences = "f"
	houses = "h"
	water  = "w"
	dirt   = "t"
)

var spriteNameMap = map[string]string{
	grass:  "res/Sprout Lands - Sprites - premium pack/tilesets/Grass.png",
	hills:  "res/Sprout Lands - Sprites - premium pack/tilesets/Hills.png",
	fences: "res/Sprout Lands - Sprites - premium pack/tilesets/Building parts/Fences.png",
	houses: "res/Sprout Lands - Sprites - premium pack/tilesets/Building parts/Wooden House.png",
	water:  "res/Sprout Lands - Sprites - premium pack/tilesets/Water.png",
	dirt:   "res/Sprout Lands - Sprites - premium pack/tilesets/Tilled Dirt.png",
}

var (
	defaultSrc    = rl.NewRectangle(0, 0, 16, 16)
	defaultDest   = defaultSrc
	defaultVector = rl.NewVector2(16, 16)
	spriteMap     = make(map[string]*rl.Texture2D)
)

type tile struct {
	texture *rl.Texture2D
	src     rl.Rectangle
}

type tiles []tile

type tilesRow []tiles

type Flooring struct {
	rows          []tilesRow
	width, height int
}

func (f Flooring) Input() {}

func (f Flooring) Update(tick int) {}

func (f Flooring) Draw() {
	for i, line := range f.rows {
		for j, tiles := range line {
			for _, tile := range tiles {
				dest := defaultDest
				dest.X = dest.Width * float32(j)
				dest.Y = dest.Height * float32(i)
				rl.DrawTexturePro(*tile.texture, tile.src, dest, defaultVector, 0, rl.White)
			}
		}
	}
}

func (f Flooring) Quit() {
	for _, line := range f.rows {
		for _, tiles := range line {
			for _, tile := range tiles {
				rl.UnloadTexture(*tile.texture)
			}
		}
	}
}

func LoadFlooring(mapFile string) (Flooring, error) {
	f, err := ioutil.ReadFile(mapFile)
	if err != nil {
		return Flooring{}, fmt.Errorf("unable to load background file. filename; %s. err: %w", mapFile, err)
	}

	lines := strings.Split(string(f), "\n")
	if len(lines) < 1 {
		return Flooring{}, fmt.Errorf("invalid background file. filename: %s", mapFile)
	}

	sizes, err := getNums(lines[0])
	if err != nil || len(sizes) != 2 {
		return Flooring{}, fmt.Errorf("invalid sizes in background file. sizes: %s", lines[0])
	}

	m := Flooring{
		rows:   []tilesRow{},
		width:  int(sizes[0]),
		height: int(sizes[1]),
	}

	if len(lines) != m.height*2+1 {
		return Flooring{}, fmt.Errorf("invalid background file. expected rows: %d, actual rows: %d", m.height*2+1, len(lines))
	}

	m.rows, err = getLines(lines, m.width)
	if err != nil {
		return Flooring{}, fmt.Errorf("invalid sizes in background file. sizes: %s", lines[0])
	}

	return m, nil
}

func getLines(lines []string, width int) ([]tilesRow, error) {
	result := []tilesRow{}

	length := (len(lines) - 1) / 2
	for i := 1; i <= length; i++ {
		nums, err := getNums(lines[i])
		if err != nil {
			return nil, fmt.Errorf("invalid nums row: %s", lines[i])
		}
		sprites, err := getSprites(lines[i+length])
		if err != nil {
			return nil, fmt.Errorf("invalid sprites row: %s", lines[i+length])
		}

		if len(nums) != len(sprites) || len(nums) != width {
			return nil, fmt.Errorf("mismatch in nums and sprites. i: %d, nums: %s (%d), sprites: %s (%d), expected length: %d", i, lines[i], len(nums), lines[i+length], len(sprites), width)
		}

		fl := tilesRow{}
		for i, num := range nums {
			tiles, err := getTiles(num, sprites[i])
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve texture. err: %w", err)
			}

			fl = append(fl, tiles)
		}

		result = append(result, fl)
	}

	return result, nil
}

func getTiles(num int32, spriteKey string) (tiles, error) {
	if num == 0 {
		return tiles{}, nil
	}

	t, err := getTexture(spriteKey)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve tiles. err: %w", err)
	}

	result := tiles{}
	if spriteKey == houses || spriteKey == fences {
		grassSprite, err := getTexture(grass)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve grass tile. err: %w", err)
		}
		result = append(result, tile{texture: grassSprite, src: defaultSrc})
	}

	width := t.Width

	src := defaultSrc
	src.X = float32((num-1)%width) * src.Width
	src.Y = float32((num-1)/width) * src.Width

	result = append(result, tile{texture: t, src: src})

	return result, nil
}

func getTexture(spriteKey string) (*rl.Texture2D, error) {
	if s, ok := spriteMap[spriteKey]; ok {
		return s, nil
	}

	if filename, ok := spriteNameMap[spriteKey]; ok {
		t := rl.LoadTexture(filename)
		spriteMap[spriteKey] = &t

		return &t, nil
	}

	return nil, fmt.Errorf("sprite not found: %s", spriteKey)
}

func getNums(line string) ([]int32, error) {
	var result []int32
	for _, col := range strings.Split(line, " ") {
		num, err := strconv.ParseInt(col, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number found. number: %s", col)
		}

		result = append(result, int32(num))
	}

	return result, nil
}

func getSprites(line string) ([]string, error) {
	var result []string
	for _, col := range strings.Split(line, " ") {
		if _, ok := spriteNameMap[col]; !ok {
			return nil, fmt.Errorf("invalid tile: %s", col)
		}

		result = append(result, col)
	}

	return result, nil
}
