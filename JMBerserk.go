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

func main() {
	pixelgl.Run(run)
}

func run() {

	win, _ := initializeWindow()

	tileList := makeTiles(win)

	floorWallSheet, _ := loadPicture("floor_wall_sheet.png")

	batch := makeWallFloorBatch(floorWallSheet)
	imd := imdraw.New(nil)

	wallFloorFrames := makeWallFloorSpriteFrames(floorWallSheet)

	batch.Clear()
	// floorBlock := pixel.NewSprite(floorWallSheet, wallFloorFrames[0])
	wallBlock := pixel.NewSprite(floorWallSheet, wallFloorFrames[1])

	for i := 0; i < len(tileList); i++ {
		if i%2 == 0 {
			wallBlock.Draw(batch, pixel.IM.Moved(tileList[i].Center()))
		} else {
			//	wallBlock.Draw(batch, pixel.IM.Moved(tileList[i].Center()))
		}
	}

	// make imd to and fill with tile rectangles
	for i := 0; i < len(tileList); i++ {
		imd.Color = colornames.White
		imd.Push(tileList[i].Min, tileList[i].Max)
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
		Title: "Wizard Berserk",
		Bounds: pixel.R(0, 0, 767, 639),
		VSync: true,
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

func makeWallFloorBatch(pic pixel.Picture) *pixel.Batch {
	batch := pixel.NewBatch(&pixel.TrianglesData{}, pic)
	return batch
}

func makeWallFloorSpriteFrames(pic pixel.Picture) []pixel.Rect {
	var floorWallSpriteFrames []pixel.Rect
	for x := pic.Bounds().Min.X; x < pic.Bounds().Max.X; x += 32 {
		for y := pic.Bounds().Min.Y; y < pic.Bounds().Max.Y; y += 32 {
			floorWallSpriteFrames = append(floorWallSpriteFrames, pixel.R(x, y, x+32, y+32))
		}
	}
	return floorWallSpriteFrames
}

// create board for level 1 - first determine which tiles will be


