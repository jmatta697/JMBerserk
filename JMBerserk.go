package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image"
	_ "image/png"
	"math/rand"
	"os"
	"time"
)

// -------------------

type hero struct {
	sprite      *pixel.Sprite
	hitBox      pixel.Rect
	lives       int
	lastDirCode float64
}

func buildHero(sprt *pixel.Sprite, freeBlockList []pixel.Rect, livesRemaining int) *hero {
	// build hero
	n := generateRandNum(len(freeBlockList))
	placementTileMin := freeBlockList[n].Min
	placementTileMax := freeBlockList[n].Max
	// make enemy's hit box location and size
	hitBoxMinScale := pixel.Vec{5, -5}
	hitBoxMaxScale := pixel.Vec{-5, 3}
	heroHitBox := pixel.Rect{placementTileMin.Add(hitBoxMinScale), placementTileMax.Add(hitBoxMaxScale)}
	heroObj := &hero{sprt, heroHitBox, livesRemaining, 0}

	return heroObj
}

// uses current hit box of heroObj to determine new position of hero hit box - returns rect obj
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

	// code below is to view hit box overlay (DEBUG)
	imd := imdraw.New(nil)
	imd.Color = colornames.White
	imd.Push(heroObj.hitBox.Min, heroObj.hitBox.Max)
	imd.Rectangle(1)
	imd.Draw(window)
}

// ----------------------

type darkMage struct {
	sprite pixel.Sprite
	hitBox pixel.Rect
}

func (enemyObj darkMage) updateEnemyHitBox(moveDirection []float64) pixel.Rect {
	deltaX := moveDirection[0]
	deltaY := moveDirection[1]
	deltaVec := pixel.Vec{float64(deltaX), float64(deltaY)}
	// gets the current hit box position and stats
	currentHitBoxMin := enemyObj.hitBox.Min
	currentHitBoxMax := enemyObj.hitBox.Max
	// set the new hit box position and stats according to the incoming move direction
	return pixel.Rect{currentHitBoxMin.Add(deltaVec), currentHitBoxMax.Add(deltaVec)}
}

// ----------------------

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

type shot struct {
	sprite    *pixel.Sprite
	hitBox    pixel.Rect
	deltaMove []float64
	active    bool
}

func buildShot(sprt *pixel.Sprite, hitBox pixel.Rect, deltaMove []float64, active bool) *shot {
	// build hero
	shotObj := &shot{sprt, hitBox, deltaMove, active}
	return shotObj
}

// determine move direction
func getShotMoveDirection(directionCode float64) []float64 {
	fmt.Println(directionCode)
	deltaPoint := []float64{0, 0}
	deltaX := 0.0
	deltaY := 0.0
	switch directionCode {
	case float64(0):
		deltaX = -0.75
	case float64(1):
		deltaX = 0.75
	case float64(2):
		deltaY = 0.75
	case float64(3):
		deltaY = -0.75
	case float64(4):
		deltaX = 0.75
		deltaY = -0.75
	case float64(5):
		deltaX = -0.75
		deltaY = -0.75
	case float64(6):
		deltaX = 0.75
		deltaY = 0.75
	case float64(7):
		deltaX = -0.75
		deltaY = 0.75
	}
	deltaPoint[0] = deltaX
	deltaPoint[1] = deltaY
	return deltaPoint
}

// uses current hit box of shotObj to determine new position of shot hit box - returns rect obj
func (shotObj shot) updateShotHitBox(deltaPoint []float64) pixel.Rect {
	deltaVec := pixel.Vec{float64(deltaPoint[0]), float64(deltaPoint[1])}
	// gets the current hit box position and stats
	currentHitBoxMin := shotObj.hitBox.Min
	currentHitBoxMax := shotObj.hitBox.Max
	// fmt.Println(currentHitBoxMin)
	// fmt.Println(currentHitBoxMax)
	// fmt.Println(pixel.Rect{currentHitBoxMin.Add(deltaVec), currentHitBoxMax.Add(deltaVec)})
	// set the new hit box position and stats according to the incoming move direction
	return pixel.Rect{currentHitBoxMin.Add(deltaVec), currentHitBoxMax.Add(deltaVec)}
}

func (shotObj shot) drawShot(window *pixelgl.Window) {
	shotMatrix := pixel.IM.Moved(shotObj.hitBox.Center())
	// currentEnemyMatrix = currentEnemyMatrix.ScaledXY(currentEnemyHitBox.Center(), pixel.Vec{10,10})
	shotObj.sprite.Draw(window, shotMatrix)

	// code below is to view hit box overlay (DEBUG)
	imd := imdraw.New(nil)
	imd.Color = colornames.White
	imd.Push(shotObj.hitBox.Min, shotObj.hitBox.Max)
	imd.Rectangle(1)
	imd.Draw(window)
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

	// make imd to and fill with tile rectangles *** DEBUG ***
	for i := 0; i < len(windowTileList); i++ {
		imd.Color = colornames.White
		imd.Push(windowTileList[i].Min, windowTileList[i].Max)
		imd.Rectangle(1)
	}

	// loaded pics from assets
	wallBlockPic, _ := loadPicture("wall_block.png")
	heroPic, _ := loadPicture("mage_0.png")
	darkMagePic, _ := loadPicture("dark_mage.png")
	shotPic, _ := loadPicture("shot.png")

	// general use wall batch
	// wallBatch := makeWallBatch(wallBlockPic)
	// wallBatch.Clear()
	// floorBlock := pixel.NewSprite(floorWallSheet, wallFloorFrames[0])

	// SPRITES
	wallBlockSprite := pixel.NewSprite(wallBlockPic, wallBlockPic.Bounds())
	heroSprite := pixel.NewSprite(heroPic, heroPic.Bounds())
	darkMageSprite := pixel.NewSprite(darkMagePic, darkMagePic.Bounds())
	shotSprite := pixel.NewSprite(shotPic, shotPic.Bounds())

	gameOver := false
	// level 1 wall setup
	level1Board := makeLevel1(wallBlockPic, wallBlockSprite, windowTileList)
	// store remaining lives
	heroLivesRemaining := 3
	// initialize various objects
	playArea := makePlayArea(level1Board, *darkMageSprite, windowTileList)
	// build hero obj
	hero := buildHero(heroSprite, playArea.freeBlockList, heroLivesRemaining)
	// build hero shot
	heroShot := buildShot(shotSprite,
		pixel.Rect{pixel.V(0, 0), pixel.V(5, 5)}, []float64{0.25, 0},
		false)
	enemyShot := buildShot(shotSprite,
		pixel.Rect{pixel.V(0, 0), pixel.V(5, 5)}, []float64{0, 0},
		false)

	secondTicker := time.NewTicker(time.Second)
	go func() {
		for range secondTicker.C {
			// Random enemy fires shot if below level-specific threshold
			if len(playArea.enemyList) != 0 && enemyShot.active == false {
				randShotChance := rand.Intn(100)
				if randShotChance < (20 * playArea.levelEnvironment.levelNum) {
					// pick random dark mage from list
					randDarkMageIndex := rand.Intn(len(playArea.enemyList))
					newEnemyShotLocationMin := playArea.enemyList[randDarkMageIndex].hitBox.Center()
					newEnemyShotLocationMax := playArea.enemyList[randDarkMageIndex].hitBox.Center().Add(pixel.V(5, 5))
					syncedEnemyShotHitBoxPosition := pixel.Rect{newEnemyShotLocationMin, newEnemyShotLocationMax}
					enemyShot.hitBox = syncedEnemyShotHitBoxPosition
					enemyShotDirectionVec := determineEnemyShotDeltaMove(hero, enemyShot.hitBox.Center())
					enemyShot.deltaMove = []float64{enemyShotDirectionVec[0], enemyShotDirectionVec[1]}
					enemyShot.active = true
				}
			}
		}
	}()

	// update all enemy hit box positions based on location of hero
	enemyMoveTicker := time.NewTicker(time.Millisecond * 50)
	// var indexRemovalList []int
	go func() {
		for range enemyMoveTicker.C {

			if len(playArea.enemyList) != 0 {
				randDarkMageIndex := generateRandNum(len(playArea.enemyList))
				enemyMoveVecDimension := determineEnemyShotDeltaMove(hero, playArea.enemyList[randDarkMageIndex].hitBox.Center())
				addVec := pixel.Vec{enemyMoveVecDimension[0], enemyMoveVecDimension[1]}
				currentEnemyHitBoxMin := playArea.enemyList[randDarkMageIndex].hitBox.Min
				currentEnemyHitBoxMax := playArea.enemyList[randDarkMageIndex].hitBox.Max
				newEnemyHitBoxMin := currentEnemyHitBoxMin.Add(addVec)
				newEnemyHitBoxMax := currentEnemyHitBoxMax.Add(addVec)
				playArea.enemyList[randDarkMageIndex].hitBox = pixel.Rect{newEnemyHitBoxMin, newEnemyHitBoxMax}

				// check collisions with walls and
				for i := len(playArea.enemyList) - 1; i >= 0; i-- {
					for j := range playArea.levelEnvironment.wallTileList {
						if len(playArea.enemyList) != 0 &&
							playArea.enemyList[i].hitBox.Intersect(playArea.levelEnvironment.wallTileList[j]) !=
								pixel.R(0, 0, 0, 0) {
							playArea.enemyList = removeEnemyFromEnemyList(playArea.enemyList, i)
							break
						}
					}
				}
			}
		}
	}()

	// main game loop
	for !win.Closed() {

		for !win.Closed() && !gameOver {
			win.Clear(colornames.Darkgrey)

			// set up cases for other levels...

			// draw all pics in level wall batch
			playArea.levelEnvironment.wallBatch.Draw(win)
			playArea.drawEnemies(win)

			// get direction of hero move from keyboard input
			heroPositionChange := checkForKeyboardInput(win)
			// fmt.Println(heroPositionChange)

			// detect a change in hero direction and set hero last position
			if heroPositionChange[0] != 0 || heroPositionChange[1] != 0 {
				// assign new last hero direction
				hero.lastDirCode = heroPositionChange[2]
			}
			// fmt.Printf("%v\n", heroPositionChange)
			// fmt.Println(hero.lastDirCode)
			// update hero hit box
			newHeroHitBox := hero.updateHitBox(heroPositionChange)
			hero.hitBox = newHeroHitBox

			// draw hero
			hero.drawHero(win)

			// check for space bar press and build shot at hero location
			if heroShot.active == false && win.Pressed(pixelgl.KeySpace) {
				// set initial location of shot hit box at hero
				newHeroShotHBMin := hero.hitBox.Center()
				newHeroShotHBMax := hero.hitBox.Center().Add(pixel.V(5, 5))
				syncedShotHitBoxPosition := pixel.Rect{newHeroShotHBMin, newHeroShotHBMax}
				heroShot.hitBox = syncedShotHitBoxPosition

				shotMoveDirection := getShotMoveDirection(hero.lastDirCode)
				fmt.Println(shotMoveDirection)

				heroShot.deltaMove = shotMoveDirection
				fmt.Println(heroShot.deltaMove)

				heroShot.active = true
			}

			// keep redrawing shot is still active
			if heroShot.active == true {
				deltaDirection := heroShot.deltaMove
				// fmt.Println(deltaDirection)
				newShotHitBox := heroShot.updateShotHitBox(deltaDirection)
				heroShot.hitBox = newShotHitBox
				heroShot.drawShot(win)
			}
			// keep redrawing enemyShot is still active
			if enemyShot.active == true {
				deltaDirection := enemyShot.deltaMove
				// fmt.Println(deltaDirection)
				newShotHitBox := enemyShot.updateShotHitBox(deltaDirection)
				enemyShot.hitBox = newShotHitBox
				enemyShot.drawShot(win)
			}

			// ------------------ collision detection ---------------------------------------

			// check for hero collision with board bounds
			if hero.hitBox.Center().Y > win.Bounds().Max.Y {
				hero.hitBox = pixel.Rect{hero.hitBox.Min.Add(pixel.Vec{0, -0.75}),
					hero.hitBox.Max.Add(pixel.Vec{0, -0.75})}
			}

			// check for hero and enemy shot collision with wallHB
			for i := 0; i < len(playArea.levelEnvironment.wallTileList); i++ {
				if heroShot.hitBox.Intersect(playArea.levelEnvironment.wallTileList[i]) !=
					pixel.R(0, 0, 0, 0) {
					heroShot.active = false
				}
				if heroShot.hitBox.Center().Y > win.Bounds().Max.Y {
					heroShot.active = false
				}
				if enemyShot.hitBox.Intersect(playArea.levelEnvironment.wallTileList[i]) !=
					pixel.R(0, 0, 0, 0) {
					enemyShot.active = false
				}
				if enemyShot.hitBox.Center().Y > win.Bounds().Max.Y {
					enemyShot.active = false
				}
			}

			// check for hero shot collision with dark mages
			for i := 0; i < len(playArea.enemyList); i++ {
				if heroShot.hitBox.Intersect(playArea.enemyList[i].hitBox) !=
					pixel.R(0, 0, 0, 0) {
					playArea.enemyList = removeEnemyFromEnemyList(playArea.enemyList, i)
					heroShot.active = false
				}
			}

			// check for hero collision with wall
			for i := 0; i < len(playArea.levelEnvironment.wallTileList); i++ {
				if hero.hitBox.Intersect(playArea.levelEnvironment.wallTileList[i]) !=
					pixel.R(0, 0, 0, 0) {
					heroLivesRemaining -= 1
					hero = buildHero(heroSprite, playArea.freeBlockList, heroLivesRemaining)
				}
			}

			// check for enemy shot with hero
			if enemyShot.hitBox.Intersect(hero.hitBox) !=
				pixel.R(0, 0, 0, 0) {
				heroLivesRemaining -= 1
				hero = buildHero(heroSprite, playArea.freeBlockList, heroLivesRemaining)
				enemyShot.active = false
			}

			if heroLivesRemaining == 0 {
				gameOver = true

			}

			// draw all tile rectangles in imd on window (for debug use)
			// imd.Draw(win)

			win.Update()
		}

		// game over screen
		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(win.Bounds().Center().Add(pixel.Vec{-120, 0}), basicAtlas)

		fmt.Fprintln(basicTxt, "GAME OVER")

		win.Clear(colornames.Black)

		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 4))
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

func removeEnemyFromEnemyList(enemyList []darkMage, indx int) []darkMage {
	return append(enemyList[:indx], enemyList[indx+1:]...)
}

// detects keyboard input and returns an array of ints [dX, dY, directionCode]
func checkForKeyboardInput(window *pixelgl.Window) []float64 {
	deltaFactor := 0.3
	deltaX := 0.0
	deltaY := 0.0

	directionCode := 0.0
	deltaPoint := []float64{0.0, 0.0, 0.0}

	if window.Pressed(pixelgl.KeyLeft) {
		deltaX -= deltaFactor
		directionCode = 0.0
	}
	if window.Pressed(pixelgl.KeyRight) {
		deltaX += deltaFactor
		directionCode = 1.0
	}
	if window.Pressed(pixelgl.KeyUp) {
		deltaY += deltaFactor
		directionCode = 2.0
	}
	if window.Pressed(pixelgl.KeyDown) {
		deltaY -= deltaFactor
		directionCode = 3.0
	}
	if window.Pressed(pixelgl.KeyDown) && window.Pressed(pixelgl.KeyRight) {
		deltaY -= deltaFactor / 100
		deltaX += deltaFactor / 100
		directionCode = 4.0
	}
	if window.Pressed(pixelgl.KeyDown) && window.Pressed(pixelgl.KeyLeft) {
		deltaY -= deltaFactor / 100
		deltaX -= deltaFactor / 100
		directionCode = 5.0
	}
	if window.Pressed(pixelgl.KeyUp) && window.Pressed(pixelgl.KeyRight) {
		deltaY += deltaFactor / 100
		deltaX += deltaFactor / 100
		directionCode = 6.0
	}
	if window.Pressed(pixelgl.KeyUp) && window.Pressed(pixelgl.KeyLeft) {
		deltaY += deltaFactor / 100
		deltaX -= deltaFactor / 100
		directionCode = 7.0
	}
	deltaPoint[0] = deltaX
	deltaPoint[1] = deltaY
	deltaPoint[2] = directionCode

	return deltaPoint
}

func determineEnemyShotDeltaMove(heroObj *hero, initialEnemyShotLocation pixel.Vec) []float64 {
	// hero center vector
	heroCenterVec := heroObj.hitBox.Center()
	differenceVec := heroCenterVec.Sub(initialEnemyShotLocation)
	unitVec := differenceVec.Unit()
	return []float64{unitVec.X, unitVec.Y}
}

// create board for level 1 - first determine which tiles will be
