package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"gogame/pkg/game"
)

const (
	screenWidth  = 1000
	screenHeight = 480
	title        = "RayGoGame"
	fps          = 60
)

func main() {
	app := &cli.App{
		Name: "gogame",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-music",
				Value: false,
				Usage: "start the game without background music",
			},
			&cli.Float64Flag{
				Aliases: []string{"s"},
				Name:    "speed",
				Value:   1.4,
				Usage:   "player speed",
			},
			&cli.Float64Flag{
				Aliases: []string{"z"},
				Name:    "zoom",
				Value:   2,
			},
			&cli.StringFlag{
				Aliases: []string{"mf"},
				Name:    "map-file",
				Value:   "res/one.map",
			},
		},
		Action: func(ctx *cli.Context) error {
			g := game.New(screenWidth, screenHeight, fps, title)
			err := g.Init(
				float32(ctx.Float64("speed")),
				ctx.Bool("no-music"),
				float32(ctx.Float64("zoom")),
				ctx.String("map-file"),
			)
			if err != nil {
				return err
			}
			defer g.Quit()

			g.Run()

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
