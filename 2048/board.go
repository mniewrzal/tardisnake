
package twenty48

import (
	"errors"
	"image/color"
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
	body []*Tile
	directionX int
	directionY int
}

// Board represents the game board.
type Board struct {
	size int
	snake  Snake
}

// NewBoard generates a new Board with giving a size.
func NewBoard(size int) (*Board, error) {
	b := &Board{
		size: size,
		snake : Snake{
			body : []*Tile {
				&Tile{
					x:0,
					y:0,
				},

			},
			directionX : 1,
			directionY : 1,
		},
	}
	go func ()  {
		time.Sleep(2 * time.Second)
		b.snake.body[0].x += b.snake.directionX
		b.snake.body[0].y += b.snake.directionY

	}()
	return b, nil
}

// Update updates the board state.
func (b *Board) Update(input *Input) error {
	// 	for t := range b.tiles {
	// 		if err := t.Update(); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	if 0 < len(b.tasks) {
	// 		t := b.tasks[0]
	// 		if err := t(); err == taskTerminated {
	// 			b.tasks = b.tasks[1:]
	// 		} else if err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	}
	// 	if dir, ok := input.Dir(); ok {
	// 		if err := b.Move(dir); err != nil {
	// 			return err
	// 		}
	// 	}
	return nil
}

// // Move enqueues tile moving tasks.
// func (b *Board) Move(dir Dir) error {
// 	for t := range b.tiles {
// 		t.stopAnimation()
// 	}
// 	if !MoveTiles(b.tiles, b.size, dir) {
// 		return nil
// 	}
// 	b.tasks = append(b.tasks, func() error {
// 		for t := range b.tiles {
// 			if t.IsMoving() {
// 				return nil
// 			}
// 		}
// 		return taskTerminated
// 	})
// 	b.tasks = append(b.tasks, func() error {
// 		nextTiles := map[*Tile]struct{}{}
// 		for t := range b.tiles {
// 			if t.IsMoving() {
// 				panic("not reach")
// 			}
// 			if t.next.value != 0 {
// 				panic("not reach")
// 			}
// 			if t.current.value == 0 {
// 				continue
// 			}
// 			nextTiles[t] = struct{}{}
// 		}
// 		b.tiles = nextTiles
// 		if err := addRandomTile(b.tiles, b.size); err != nil {
// 			return err
// 		}
// 		return taskTerminated
// 	})
// 	return nil
// }

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
		y :=  tile.y*tileSize + ( tile.y+1)*tileMargin
		op.GeoM.Translate(float64(x), float64(y))
	
		r, g, b, a := colorToScale(color.NRGBA{0xee, 0xFF, 0xFF, 0xFF})
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
