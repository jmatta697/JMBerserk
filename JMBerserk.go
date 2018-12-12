package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"math/rand"
	"os"
)

type hero struct {
	hitBox pixel.Rect
	sprite pixel.Sprite
	lives  int
}

type darkMage struct {
	sprite pixel.Sprite
	hitBox pixel.Rect
}

type level struct {
	wallTileList []pixel.Rect
	wallBatch    *pixel.Batch
	levelDescrpt string
	levelNum     int
}

type playArea struct {
	levelEnvironment *level
	enemyList        []darkMage
}

// --------------------------------

func main() {
	pixelgl.Run(run)
}

func run() {

	win, _ := initializeWindow()

	// all window tiles
	windowTileList := makeTiles(win)

	// loaded pics from assets
	wallBlockPic, _ := loadPicture("wall_block.png")
	heroPic, _ := loadPicture("mage_0.png")
	darkMagePic, _ := loadPicture("dark_mage.png")

	// general use wall batch
	// wallBatch := makeWallBatch(wallBlockPic)

	imd := imdraw.New(nil)

	// wallBatch.Clear()
	// floorBlock := pixel.NewSprite(floorWallSheet, wallFloorFrames[0])
	wallBlockSprite := pixel.NewSprite(wallBlockPic, wallBlockPic.Bounds())
	heroSprite := pixel.NewSprite(heroPic, heroPic.Bounds())
	darkMageSprite := pixel.NewSprite(darkMagePic, darkMagePic.Bounds())

	// level 1 wall setup
	level1Board := makeLevel1(wallBlockPic, wallBlockSprite, windowTileList)

	// make imd to and fill with tile rectangles
	for i := 0; i < len(windowTileList); i++ {
		imd.Color = colornames.White
		imd.Push(windowTileList[i].Min, windowTileList[i].Max)
		imd.Rectangle(1)
	}

	// main game loop
	for !win.Closed() {

		// gameOver := false

		for !win.Closed() {
			win.Clear(colornames.Darkgrey)

			// set up cases for other levels...

			// draw all pics in level wall batch
			level1Board.wallBatch.Draw(win)

			// draw all tile rectangles in imd on window (for debug use)
			imd.Draw(win)

			win.Update()
		}

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

// builds level1 object
func makeLevel1(wallPic pixel.Picture, wallSprite *pixel.Sprite, winTileLst []pixel.Rect) level {

	wallBatch := makeWallBatch(wallPic)
	var level1WallList []pixel.Rect

	// wallList holds all rects
	level1 := level{level1WallList, wallBatch, "Level 1", 1}

	// Add wall tiles appropriate to level1 and draw wall sprites to level1 batch
	for i := 0; i < 24; i++ {
		level1.wallTileList = append(level1.wallTileList, winTileLst[i])
		wallSprite.Draw(wallBatch, pixel.IM.Moved(winTileLst[i].Center()))
	}
	for i := 24; i < 456; i += 24 {
		level1.wallTileList = append(level1.wallTileList, winTileLst[i])
		wallSprite.Draw(wallBatch, pixel.IM.Moved(winTileLst[i].Center()))
	}
	for i := 47; i < 456; i += 24 {
		level1.wallTileList = append(level1.wallTileList, winTileLst[i])
		wallSprite.Draw(wallBatch, pixel.IM.Moved(winTileLst[i].Center()))
	}
	for i := 456; i < 465; i++ {
		level1.wallTileList = append(level1.wallTileList, winTileLst[i])
		wallSprite.Draw(wallBatch, pixel.IM.Moved(winTileLst[i].Center()))
	}
	for i := 471; i < len(winTileLst); i++ {
		level1.wallTileList = append(level1.wallTileList, winTileLst[i])
		wallSprite.Draw(wallBatch, pixel.IM.Moved(winTileLst[i].Center()))
	}
	for i := 125; i < 342; i += 24 {
		level1.wallTileList = append(level1.wallTileList, winTileLst[i])
		wallSprite.Draw(wallBatch, pixel.IM.Moved(winTileLst[i].Center()))
	}
	for i := 246; i < 258; i++ {
		level1.wallTileList = append(level1.wallTileList, winTileLst[i])
		wallSprite.Draw(wallBatch, pixel.IM.Moved(winTileLst[i].Center()))
	}
	for i := 138; i < 355; i += 24 {
		level1.wallTileList = append(level1.wallTileList, winTileLst[i])
		wallSprite.Draw(wallBatch, pixel.IM.Moved(winTileLst[i].Center()))
	}

	return level1
}

func makeEnemyBatch(pic pixel.Picture) *pixel.Batch {
	batch := pixel.NewBatch(&pixel.TrianglesData{}, pic)
	return batch
}

func makePlayArea(lvl *level, hroObj *hero, enemySpite pixel.Sprite, winTileList []pixel.Rect) *playArea {
	var playArea playArea

	// initialize play area enemy list
	var enemyList []darkMage

	// this is all unoccupied blocks in this play area (all window tiles - level wall blocks)
	availableBlockList := getAvailableTiles(winTileList, lvl.wallTileList)

	// build enemy list
	for i := 0; i < 4+(2*lvl.levelNum); i++ {
		// random element index for picking out a block to place an enemy
		n := rand.Int() % len(availableBlockList)
		// get tile info
		placementTileMin := availableBlockList[n].Min
		placementTileMax := availableBlockList[n].Max
		// make enemy's hit box location and size
		enemyHitBox := pixel.Rect{placementTileMin, placementTileMax}
		darkMageObj := darkMage{enemySpite, enemyHitBox}
		enemyList = append(enemyList, darkMageObj)
	}

}

// subtracts occupied tiles from list off all window tiles
func getAvailableTiles(allWinTiles []pixel.Rect, occupiedList []pixel.Rect) []pixel.Rect {

	var availableBlocks []pixel.Rect
	var validSpot bool

	for i := 0; i < len(allWinTiles); i++ {
		validSpot = true
		for j := 0; j < len(occupiedList); i++ {
			if allWinTiles[i].Center() == occupiedList[j].Center() {
				validSpot = false
			}
		}
		if validSpot {
			availableBlocks = append(availableBlocks, allWinTiles[i])
		}
	}

	return availableBlocks
}

// create board for level 1 - first determine which tiles will be
