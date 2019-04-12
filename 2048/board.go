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

const (
	arcadeFontBaseSize = 8
)

var (
	backgroundColor = color.RGBA{0x17, 0x9c, 0xc1, 0xff}
	tileSize        = 50
	tileMargin      = 0

	frameColor = color.RGBA{0xbb, 0xad, 0xa0, 0xff}

	tileImage        *ebiten.Image
	arcadeFonts      map[int]font.Face
	chosenFood       *ebiten.Image
	foodImage        *ebiten.Image
	foodImage2       *ebiten.Image
	backgroundImage  *ebiten.Image
	backgroundImage2 *ebiten.Image
	backgroundImage3 *ebiten.Image
	backgroundImage4 *ebiten.Image

	tardiHeadImageUp    *ebiten.Image
	tardiHeadImageRight *ebiten.Image
	tardiHeadImageDown  *ebiten.Image
	tardiHeadImageLeft  *ebiten.Image
)

func init() {
	tileImage, _ = ebiten.NewImage(tileSize, tileSize, ebiten.FilterDefault)
	tileImage.Fill(color.White)

	foodImage, _, _ = ebitenutil.NewImageFromFile("folder.png", ebiten.FilterDefault)
	foodImage2, _, _ = ebitenutil.NewImageFromFile("bucket.png", ebiten.FilterDefault)
	var err error
	backgroundImage, _, err = ebitenutil.NewImageFromFile("space.png", ebiten.FilterDefault)
	if err != nil {
		fmt.Println(err)
	}
	backgroundImage2, _, err = ebitenutil.NewImageFromFile("desert.png", ebiten.FilterDefault)
	if err != nil {
		fmt.Println(err)
	}
	backgroundImage3, _, err = ebitenutil.NewImageFromFile("underwater.png", ebiten.FilterDefault)
	if err != nil {
		fmt.Println(err)
	}
	backgroundImage4, _, err = ebitenutil.NewImageFromFile("snowy.png", ebiten.FilterDefault)
	if err != nil {
		fmt.Println(err)
	}

	tardiHeadImageUp, _, _ = ebitenutil.NewImageFromFile("t1.png", ebiten.FilterDefault)
	tardiHeadImageRight, _, _ = ebitenutil.NewImageFromFile("t2.png", ebiten.FilterDefault)
	tardiHeadImageDown, _, _ = ebitenutil.NewImageFromFile("t3.png", ebiten.FilterDefault)
	tardiHeadImageLeft, _, _ = ebitenutil.NewImageFromFile("t4.png", ebiten.FilterDefault)
}

type Tile struct {
	x    int
	y    int
	dirX int
	dirY int
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
	sounds   *sounds
	snake    Snake
	food     Food
	playMode bool
}

// NewBoard generates a new Board with giving a size.
func NewBoard(sounds *sounds, size int) (*Board, error) {
	b := &Board{
		size:   size,
		sounds: sounds,
		snake: Snake{
			body: []*Tile{
				&Tile{
					x:    2,
					y:    0,
					dirX: 1,
					dirY: 0,
				},
				&Tile{
					x:    1,
					y:    0,
					dirX: 1,
					dirY: 0,
				}, &Tile{
					x:    0,
					y:    0,
					dirX: 1,
					dirY: 0,
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
	go func() {
		for {
			if !b.playMode {
				sounds.music.Pause()
			}
			if !sounds.music.IsPlaying() {
				sounds.music.Rewind()
				sounds.music.Play()
			}
		}
	}()
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
					sounds.death.Play()
					fmt.Println("Snake ate itself. Game over :(")
					return
				}
			}

			for i := len(b.snake.body) - 1; i > 0; i-- {
				b.snake.body[i].x = b.snake.body[i-1].x
				b.snake.body[i].y = b.snake.body[i-1].y

				b.snake.body[i].dirX = b.snake.body[i-1].dirX
				b.snake.body[i].dirY = b.snake.body[i-1].dirY
			}
			if head.x <= b.size {
				if head.x == 0 && b.snake.directionX == -1 {
					head.x = b.size + 1
				}
				head.x += b.snake.directionX
				head.dirX = b.snake.directionX
			} else {
				head.x -= b.size + 1
				head.dirX = b.snake.directionX
			}

			if head.y <= b.size {
				if head.y == 0 && b.snake.directionY == -1 {
					head.y = b.size + 1
				}
				head.y += b.snake.directionY
				head.dirY = b.snake.directionY
			} else {
				head.y -= b.size + 1
				head.dirY = b.snake.directionY
			}

			// check if food has been found
			if b.food.tiles[0] != nil {
				for _, f := range b.food.tiles {
					if head.x == f.x && head.y == f.y {
						sounds.score.Rewind()
						sounds.score.Play()

						r := rand.Intn(2)
						if r == 0 {
							chosenFood = foodImage
						} else {
							chosenFood = foodImage2
						}
						b.generateFood()
						tail := len(b.snake.body) - 1
						newTile := &Tile{
							x:    b.snake.body[tail].x,
							y:    b.snake.body[tail].y,
							dirX: b.snake.body[tail].dirX,
							dirY: b.snake.body[tail].dirY,
						}
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
	// for j := 0; j < board.size; j++ {
	// 	for i := 0; i < board.size; i++ {
	// 		// v := 0
	// 		op := &ebiten.DrawImageOptions{}
	// 		x := i * tileSize
	// 		y := j * tileSize
	// 		op.GeoM.Translate(float64(x), float64(y))

	// 		r, g, b, a := colorToScale(color.NRGBA{0xee, 0xe4, 0xda, 0x59})
	// 		op.ColorM.Scale(r, g, b, a)
	// 		boardImage.DrawImage(tileImage, op)
	// 	}
	// }
	op := &ebiten.DrawImageOptions{}

	boardImage.DrawImage(backgroundImage, op)

	if board.GetPoints() >= 3 {
		boardImage.DrawImage(backgroundImage2, op)
	}
	if board.GetPoints() >= 6 {
		boardImage.DrawImage(backgroundImage3, op)
	}
	if board.GetPoints() >= 9 {
		boardImage.DrawImage(backgroundImage4, op)
	}

	for i, tile := range board.snake.body {
		if i == 0 {
			op := &ebiten.DrawImageOptions{}
			x := tile.x * tileSize
			y := tile.y * tileSize
			op.GeoM.Translate(float64(x), float64(y))

			var headImage *ebiten.Image
			if board.snake.directionX == 1 {
				headImage = tardiHeadImageRight
			} else if board.snake.directionX == -1 {
				headImage = tardiHeadImageLeft
			} else if board.snake.directionY == 1 {
				headImage = tardiHeadImageDown
			} else if board.snake.directionY == -1 {
				headImage = tardiHeadImageUp
			}

			boardImage.DrawImage(headImage, op)
		} else {
			op := &ebiten.DrawImageOptions{}
			x := tile.x * tileSize
			y := tile.y * tileSize
			op.GeoM.Translate(float64(x), float64(y))

			var bodyImage *ebiten.Image
			if tile.dirX == 1 {
				bodyImage = tardiHeadImageRight
			} else if tile.dirX == -1 {
				bodyImage = tardiHeadImageLeft
			} else if tile.dirY == 1 {
				bodyImage = tardiHeadImageDown
			} else if tile.dirY == -1 {
				bodyImage = tardiHeadImageUp
			}

			boardImage.DrawImage(bodyImage, op)
		}

	}

	for _, tile := range board.food.tiles {
		if tile == nil {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		x := tile.x * tileSize
		y := tile.y * tileSize
		op.GeoM.Translate(float64(x), float64(y))

		if chosenFood == nil {
			chosenFood = foodImage
		}
		boardImage.DrawImage(chosenFood, op)
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
