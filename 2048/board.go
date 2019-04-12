package twenty48

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/font"
)

var (
	tileSize   = 50
	tileMargin = 4
)

var (
	backgroundColor = color.RGBA{0xfa, 0xf8, 0xef, 0xff}
	frameColor      = color.RGBA{0xbb, 0xad, 0xa0, 0xff}
)

const (
	arcadeFontBaseSize = 8
)

var (
	tileImage   *ebiten.Image
	arcadeFonts map[int]font.Face
	foodImage   *ebiten.Image
)

func init() {
	tileImage, _ = ebiten.NewImage(tileSize, tileSize, ebiten.FilterDefault)
	tileImage.Fill(color.White)

	foodImage, _, _ = ebitenutil.NewImageFromFile("encrypt.png", ebiten.FilterDefault)
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
	size     int
	snake    Snake
	food     Food
	playMode bool
}

// NewBoard generates a new Board with giving a size.
func NewBoard(size int) (*Board, error) {
	b := &Board{
		size: size,
		snake: Snake{
			body: []*Tile{
				&Tile{
					x: 2,
					y: 0,
				},
				&Tile{
					x: 1,
					y: 0,
				}, &Tile{
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
		playMode: true,
	}
	head := b.snake.body[0]
	go func() {
		for {
			time.Sleep(150 * time.Millisecond)

			next := &Tile{
				x: head.x + b.snake.directionX,
				y: head.y + b.snake.directionY,
			}
			if exists(b.snake.body, next) {
				if exists([]*Tile{next}, b.snake.body[1]) {
					if b.snake.directionX != 0 {
						b.snake.directionX *= -1
					}
					if b.snake.directionY != 0 {
						b.snake.directionY *= -1
					}
				} else {
					b.playMode = false
					fmt.Println("Snake ate itself. Game over :(")
					return
				}
			}

			for i := len(b.snake.body) - 1; i > 0; i-- {
				b.snake.body[i].x = b.snake.body[i-1].x
				b.snake.body[i].y = b.snake.body[i-1].y
			}
			if head.x <= b.size {
				head.x += b.snake.directionX
			} else {
				head.x -= b.size + 1
			}

			if head.y <= b.size {
				head.y += b.snake.directionY
			} else {
				head.y -= b.size + 1
			}

			// check if food has been found
			if b.food.tiles[0] != nil {
				for _, f := range b.food.tiles {
					if head.x == f.x && head.y == f.y {
						b.generateFood()
						tail := len(b.snake.body) - 1
						newTile := &Tile{x: b.snake.body[tail].x, y: b.snake.body[tail].y}
						b.snake.body = append(b.snake.body, newTile)
					}
				}
			}
		}
	}()
	go func() {
		for b.playMode {
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

// GetPoints returns the points won so far
func (b *Board) GetPoints() int {
	return len(b.snake.body) - 3
}

// Draw draws the board to the given boardImage.
func (board *Board) Draw(boardImage *ebiten.Image) {
	boardImage.Fill(frameColor)
	for j := 0; j < board.size; j++ {
		for i := 0; i < board.size; i++ {
			// v := 0
			op := &ebiten.DrawImageOptions{}
			x := i * tileSize
			y := j * tileSize
			op.GeoM.Translate(float64(x), float64(y))

			r, g, b, a := colorToScale(color.NRGBA{0xee, 0xe4, 0xda, 0x59})
			op.ColorM.Scale(r, g, b, a)
			boardImage.DrawImage(tileImage, op)
		}
	}

	for _, tile := range board.snake.body {
		op := &ebiten.DrawImageOptions{}
		x := tile.x * tileSize
		y := tile.y * tileSize
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
		x := tile.x * tileSize
		y := tile.y * tileSize
		op.GeoM.Translate(float64(x), float64(y))

		r, g, b, a := colorToScale(color.NRGBA{0xee, 0xAA, 0xAA, 0xAA})
		op.ColorM.Scale(r, g, b, a)
		boardImage.DrawImage(foodImage, op)
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
