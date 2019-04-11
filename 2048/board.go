package twenty48

import (
	"errors"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

var taskTerminated = errors.New("twenty48: task terminated")

var (
	tileSize   = 50
	tileMargin = 4
)

var (
	backgroundColor = color.RGBA{0xfa, 0xf8, 0xef, 0xff}
	frameColor      = color.RGBA{0xbb, 0xad, 0xa0, 0xff}
)

var (
	tileImage *ebiten.Image
)

func init() {
	tileImage, _ = ebiten.NewImage(tileSize, tileSize, ebiten.FilterDefault)
	tileImage.Fill(color.White)
}

type task func() error

type Tile struct {
	x int
	y int
}

type Snake struct {
	body       []*Tile
	directionX int
	directionY int
}

type Food struct {
	tiles []*Tile
}

// Board represents the game board.
type Board struct {
	size  int
	snake Snake
	food  Food
}

// NewBoard generates a new Board with giving a size.
func NewBoard(size int) (*Board, error) {
	b := &Board{
		size: size,
		snake: Snake{
			body: []*Tile{
				&Tile{
					x: 0,
					y: 0,
				},
			},
			directionX: 1,
			directionY: 0,
		},
		food: Food{
			tiles: make([]*Tile, 1),
		},
	}
	go func() {
		for {
			time.Sleep(1 * time.Second)
			b.snake.body[0].x += b.snake.directionX
			b.snake.body[0].y += b.snake.directionY
		}
	}()
	go func() {
		for {
			time.Sleep(4 * time.Second)
			b.generateFood()
		}
	}()
	return b, nil
}

func (b *Board) generateFood() {
	for {
		newFoodTile := &Tile{
			x: rand.Intn(b.size),
			y: rand.Intn(b.size),
		}

		if !exists(b.snake.body, newFoodTile) {
			b.food.tiles[0] = newFoodTile
			return
		}
	}
}

func exists(tiles []*Tile, target *Tile) bool {
	for _, tile := range tiles {
		if target.x == tile.x && target.y == tile.y {
			return true
		}
	}
	return false
}

// Update updates the board state.
func (b *Board) Update(input *Input) error {
	if dir, ok := input.Dir(); ok {
		switch dir {
		case DirDown:
			b.snake.directionX = 0
			b.snake.directionY = 1
		case DirUp:
			b.snake.directionX = 0
			b.snake.directionY = -1
		case DirLeft:
			b.snake.directionX = -1
			b.snake.directionY = 0
		case DirRight:
			b.snake.directionX = 1
			b.snake.directionY = 0
		}
	}
	return nil
}

// Size returns the board size.
func (b *Board) Size() (int, int) {
	x := b.size*tileSize + (b.size+1)*tileMargin
	y := x
	return x, y
}

// Draw draws the board to the given boardImage.
func (board *Board) Draw(boardImage *ebiten.Image) {
	boardImage.Fill(frameColor)
	for j := 0; j < board.size; j++ {
		for i := 0; i < board.size; i++ {
			// v := 0
			op := &ebiten.DrawImageOptions{}
			x := i*tileSize + (i+1)*tileMargin
			y := j*tileSize + (j+1)*tileMargin
			op.GeoM.Translate(float64(x), float64(y))

			r, g, b, a := colorToScale(color.NRGBA{0xee, 0xe4, 0xda, 0x59})
			op.ColorM.Scale(r, g, b, a)
			boardImage.DrawImage(tileImage, op)
		}
	}

	for _, tile := range board.snake.body {
		op := &ebiten.DrawImageOptions{}
		x := tile.x*tileSize + (tile.x+1)*tileMargin
		y := tile.y*tileSize + (tile.y+1)*tileMargin
		op.GeoM.Translate(float64(x), float64(y))

		r, g, b, a := colorToScale(color.NRGBA{0xee, 0xFF, 0xFF, 0xFF})
		op.ColorM.Scale(r, g, b, a)
		boardImage.DrawImage(tileImage, op)
	}

	for _, tile := range board.food.tiles {
		if tile == nil {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		x := tile.x*tileSize + (tile.x+1)*tileMargin
		y := tile.y*tileSize + (tile.y+1)*tileMargin
		op.GeoM.Translate(float64(x), float64(y))

		r, g, b, a := colorToScale(color.NRGBA{0xee, 0xAA, 0xAA, 0xAA})
		op.ColorM.Scale(r, g, b, a)
		boardImage.DrawImage(tileImage, op)
	}

}

func colorToScale(clr color.Color) (float64, float64, float64, float64) {
	r, g, b, a := clr.RGBA()
	rf := float64(r) / 0xffff
	gf := float64(g) / 0xffff
	bf := float64(b) / 0xffff
	af := float64(a) / 0xffff
	// Convert to non-premultiplied alpha components.
	if 0 < af {
		rf /= af
		gf /= af
		bf /= af
	}
	return rf, gf, bf, af
}
