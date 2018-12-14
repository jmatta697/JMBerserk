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
	"time"
)

type hero struct {
	sprite            *pixel.Sprite
	hitBox            pixel.Rect
	lives             int
	lastDirectionMove []int
}

func buildHero(sprt *pixel.Sprite, freeBlockList []pixel.Rect) hero {
	// build hero
	n := generateRandNum(len(freeBlockList))
	placementTileMin := freeBlockList[n].Min
	placementTileMax := freeBlockList[n].Max
	// make enemy's hit box location and size
	hitBoxMinScale := pixel.Vec{5, -5}
	hitBoxMaxScale := pixel.Vec{-5, 3}
	heroHitBox := pixel.Rect{placementTileMin.Add(hitBoxMinScale), placementTileMax.Add(hitBoxMaxScale)}
	heroObj := hero{sprt, heroHitBox, 3, []int{0, 0}}

	return heroObj
}

func (heroObj hero) updateHitBox(moveDirection []float64) pixel.Rect {
	deltaX := moveDirection[0]
	deltaY := moveDirection[1]
	deltaVec := pixel.Vec{float64(deltaX), float64(deltaY)}
	// gets the current hit box position and stats
	currentHitBoxMin := heroObj.hitBox.Min
	currentHitBoxMax := heroObj.hitBox.Max
	// set the new hit box position and stats according to the incoming move direction

	return pixel.Rect{currentHitBoxMin.Add(deltaVec), currentHitBoxMax.Add(deltaVec)}
}

func (heroObj hero) drawHero(window *pixelgl.Window) {
	heroMatrix := pixel.IM.Moved(heroObj.hitBox.Center())
	// currentEnemyMatrix = currentEnemyMatrix.ScaledXY(currentEnemyHitBox.Center(), pixel.Vec{10,10})
	heroObj.sprite.Draw(window, heroMatrix)

	imd := imdraw.New(nil)
	imd.Color = colornames.White
	imd.Push(heroObj.hitBox.Min, heroObj.hitBox.Max)
	imd.Rectangle(1)
	imd.Draw(window)
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

// ------------------------------

type playArea struct {
	levelEnvironment level
	enemyList        []darkMage
	freeBlockList    []pixel.Rect
}

func (pa playArea) drawEnemies(window *pixelgl.Window) {
	for i := range pa.enemyList {
		currentEnemy := pa.enemyList[i]
		currentEnemyHitBox := currentEnemy.hitBox
		currentEnemyMatrix := pixel.IM
		// currentEnemyMatrix = currentEnemyMatrix.ScaledXY(currentEnemyHitBox.Center(), pixel.Vec{10,10})
		currentEnemyMatrix = currentEnemyMatrix.Moved(currentEnemyHitBox.Center())
		currentEnemy.sprite.Draw(window, currentEnemyMatrix)
	}
}

// --------------------------------

func main() {
	pixelgl.Run(run)
}

func run() {

	win, _ := initializeWindow()

	// all window tiles
	windowTileList := makeTiles(win)

	imd := imdraw.New(nil)

	// make imd to and fill with tile rectangles
	for i := 0; i < len(windowTileList); i++ {
		imd.Color = colornames.White
		imd.Push(windowTileList[i].Min, windowTileList[i].Max)
		imd.Rectangle(1)
	}

	// loaded pics from assets
	wallBlockPic, _ := loadPicture("wall_block.png")
	heroPic, _ := loadPicture("mage_0.png")
	darkMagePic, _ := loadPicture("dark_mage.png")

	// general use wall batch
	// wallBatch := makeWallBatch(wallBlockPic)
	// wallBatch.Clear()
	// floorBlock := pixel.NewSprite(floorWallSheet, wallFloorFrames[0])
	wallBlockSprite := pixel.NewSprite(wallBlockPic, wallBlockPic.Bounds())
	heroSprite := pixel.NewSprite(heroPic, heroPic.Bounds())
	darkMageSprite := pixel.NewSprite(darkMagePic, darkMagePic.Bounds())

	// level 1 wall setup
	level1Board := makeLevel1(wallBlockPic, wallBlockSprite, windowTileList)

	// main game loop
	for !win.Closed() {

		// gameOver := false
		playArea := makePlayArea(level1Board, *darkMageSprite, windowTileList)
		// build hero obj
		heroStrctObj := buildHero(heroSprite, playArea.freeBlockList)

		for !win.Closed() {
			win.Clear(colornames.Darkgrey)

			// set up cases for other levels...

			// draw all pics in level wall batch
			playArea.levelEnvironment.wallBatch.Draw(win)
			playArea.drawEnemies(win)

			// get direction of hero move from keyboard input
			heroPositionChange := checkForKeyboardInput(win)
			// update hero hit box
			newHeroHitBox := heroStrctObj.updateHitBox(heroPositionChange)
			heroStrctObj.hitBox = newHeroHitBox
			// draw hero
			heroStrctObj.drawHero(win)

			// draw all tile rectangles in imd on window (for debug use)
			// imd.Draw(win)

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

func makePlayArea(lvl level, enemySpite pixel.Sprite, winTileList []pixel.Rect) playArea {
	var playArea playArea

	// initialize play area enemy list
	var enemyList []darkMage

	// this is all unoccupied blocks in this play area (all window tiles - level wall blocks)
	availableBlockList := getAvailableTiles(winTileList, lvl.wallTileList)

	// build enemy list
	for i := 0; i < 4+(2*lvl.levelNum); i++ {
		// random element index for picking out a block to place an enemy
		n := generateRandNum(len(availableBlockList))
		// get tile info
		placementTileMin := availableBlockList[n].Min
		placementTileMax := availableBlockList[n].Max
		// make enemy's hit box location and size
		enemyHitBox := pixel.Rect{placementTileMin, placementTileMax}
		darkMageObj := darkMage{enemySpite, enemyHitBox}
		enemyList = append(enemyList, darkMageObj)
		// remove the block that is taken up by the newly added enemy
		availableBlockList = removeAvailableSpotFromList(availableBlockList, n)
	}

	playArea.levelEnvironment = lvl
	playArea.enemyList = enemyList
	playArea.freeBlockList = availableBlockList

	return playArea

}

// subtracts occupied tiles from list off all window tiles
func getAvailableTiles(allWinTiles []pixel.Rect, occupiedList []pixel.Rect) []pixel.Rect {

	var availableBlocks []pixel.Rect
	var validSpot bool

	for i := 0; i < len(allWinTiles); i++ {
		validSpot = true
		for j := 0; j < len(occupiedList); j++ {
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

// returns a random integer from 0 to max
func generateRandNum(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

// removes element at index (indx) from list (rectLst) and return new list
func removeAvailableSpotFromList(rectLst []pixel.Rect, indx int) []pixel.Rect {
	return append(rectLst[:indx], rectLst[indx+1:]...)
}

// detects keyboard input and returns an array of ints [dX, dY, directionCode]
func checkForKeyboardInput(window *pixelgl.Window) []float64 {
	deltaFactor := 0.25
	deltaX := 0.0
	deltaY := 0.0

	// directionCode := 0
	deltaPoint := []float64{0.0, 0.0}

	if window.Pressed(pixelgl.KeyLeft) {
		deltaX -= deltaFactor * 2
		// directionCode = 0
	}
	if window.Pressed(pixelgl.KeyRight) {
		deltaX += deltaFactor * 2
		// directionCode = 1
	}
	if window.Pressed(pixelgl.KeyUp) {
		deltaY += deltaFactor * 2
		// directionCode = 2
	}
	if window.Pressed(pixelgl.KeyDown) {
		deltaY -= deltaFactor * 2
		// directionCode = 3
	}
	if window.Pressed(pixelgl.KeyDown) && window.Pressed(pixelgl.KeyRight) {
		deltaY -= deltaFactor
		deltaX += deltaFactor
		// directionCode = 4
	}
	if window.Pressed(pixelgl.KeyDown) && window.Pressed(pixelgl.KeyLeft) {
		deltaY -= deltaFactor
		deltaX -= deltaFactor
		// directionCode = 5
	}
	if window.Pressed(pixelgl.KeyUp) && window.Pressed(pixelgl.KeyRight) {
		deltaY += deltaFactor
		deltaX += deltaFactor
		// directionCode = 6
	}
	if window.Pressed(pixelgl.KeyUp) && window.Pressed(pixelgl.KeyLeft) {
		deltaY += deltaFactor
		deltaX -= deltaFactor
		// directionCode = 7
	}
	deltaPoint[0] = deltaX
	deltaPoint[1] = deltaY

	return deltaPoint
}

// create board for level 1 - first determine which tiles will be
