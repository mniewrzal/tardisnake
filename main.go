// Copyright 2016 The Ebiten Authors



package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/mniewrzal/tardisnake/2048"
)

var (
	game *twenty48.Game
)

func update(screen *ebiten.Image) error {
	if err := game.Update(); err != nil {
		return err
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	game.Draw(screen)
	return nil
}

func main() {
	var err error
	game, err = twenty48.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.Run(update, twenty48.ScreenWidth, twenty48.ScreenHeight, 1, "2048 (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}
