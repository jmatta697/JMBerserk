package main

import (
	"fmt"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"testing"
)

func TestMain(m *testing.M) {
	var t *testing.T
	go TestMakeTiles(t)
}


func TestMakeTiles(t *testing.T) {
	win, _ := initializeWindow()

	imd := imdraw.New(nil)

	tilelist := makeTiles(win)

	fmt.Println(tilelist)

	for i := 0; i < len(tilelist); i++ {
		imd.Color = colornames.White
		imd.Push(tilelist[i].Min, tilelist[i].Max)
		imd.Rectangle(10)
	}

	// main game loop
	for !win.Closed() {

		win.Clear(colornames.White)

		imd.Draw(win)

		win.Update()
	}



}
