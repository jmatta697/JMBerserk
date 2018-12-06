package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"os"
)

type hero struct {
	hitBox pixel.Rect
	sprite pixel.Sprite
	lives  int
}

type darkMage struct {
	hitBox pixel.Rect
	sprite pixel.Sprite
}

type levelBoard struct {
	wallTileList []pixel.Rect
	levelDescrpt string
}

func main() {
	pixelgl.Run(run)
}

func run() {

	win, _ := initializeWindow()

	windowTileList := makeTiles(win)

	wallBlockPic, _ := loadPicture("wall_block.png")

	batch := makeWallBatch(wallBlockPic)
	imd := imdraw.New(nil)

	batch.Clear()
	// floorBlock := pixel.NewSprite(floorWallSheet, wallFloorFrames[0])
	wallBlockSprite := pixel.NewSprite(wallBlockPic, wallBlockPic.Bounds())

	// level 1 wall setup
	var level1WallList []pixel.Rect
	level1Board := levelBoard{level1WallList, "Level 1"}
	for i := 0; i < 24; i++ {
		level1Board.wallTileList = append(level1Board.wallTileList, windowTileList[i])
		wallBlockSprite.Draw(batch, pixel.IM.Moved(windowTileList[i].Center()))
	}
	for i := 24; i < 456; i += 24 {
		level1Board.wallTileList = append(level1Board.wallTileList, windowTileList[i])
		wallBlockSprite.Draw(batch, pixel.IM.Moved(windowTileList[i].Center()))
	}
	for i := 47; i < 456; i += 24 {
		level1Board.wallTileList = append(level1Board.wallTileList, windowTileList[i])
		wallBlockSprite.Draw(batch, pixel.IM.Moved(windowTileList[i].Center()))
	}
	for i := 456; i < 465; i++ {
		level1Board.wallTileList = append(level1Board.wallTileList, windowTileList[i])
		wallBlockSprite.Draw(batch, pixel.IM.Moved(windowTileList[i].Center()))
	}
	for i := 471; i < len(windowTileList); i++ {
		level1Board.wallTileList = append(level1Board.wallTileList, windowTileList[i])
		wallBlockSprite.Draw(batch, pixel.IM.Moved(windowTileList[i].Center()))
	}
	for i := 125; i < 342; i += 24 {
		level1Board.wallTileList = append(level1Board.wallTileList, windowTileList[i])
		wallBlockSprite.Draw(batch, pixel.IM.Moved(windowTileList[i].Center()))
	}
	for i := 246; i < 258; i++ {
		level1Board.wallTileList = append(level1Board.wallTileList, windowTileList[i])
		wallBlockSprite.Draw(batch, pixel.IM.Moved(windowTileList[i].Center()))
	}

	// make imd to and fill with tile rectangles
	for i := 0; i < len(windowTileList); i++ {
		imd.Color = colornames.White
		imd.Push(windowTileList[i].Min, windowTileList[i].Max)
		imd.Rectangle(1)
	}

	// main game loop
	for !win.Closed() {

		win.Clear(colornames.Darkgrey)

		// draw all pics in batch
		batch.Draw(win)
		// draw all tile rectangles in imd on window
		imd.Draw(win)

		win.Update()
	}
}

// ------------------- utility functions -----------------

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func initializeWindow() (*pixelgl.Window, pixelgl.WindowConfig) {
	cfg := pixelgl.WindowConfig{
		Title:  "Wizard Berserk",
		Bounds: pixel.R(0, 0, 767, 639),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	return win, cfg
}

func makeTiles(window *pixelgl.Window) []pixel.Rect {
	var tiles []pixel.Rect
	for y := window.Bounds().Min.Y; y < window.Bounds().Max.Y; y += 32 {
		for x := window.Bounds().Min.X; x < window.Bounds().Max.X; x += 32 {
			tiles = append(tiles, pixel.R(x, y, x+32, y+32))
		}
	}
	return tiles
}

func makeWallBatch(pic pixel.Picture) *pixel.Batch {
	batch := pixel.NewBatch(&pixel.TrianglesData{}, pic)
	return batch
}

// create board for level 1 - first determine which tiles will be
