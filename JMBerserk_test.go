package main

import (
	"github.com/faiface/pixel"
	"testing"
)

func TestGetAvailableTiles(t *testing.T) {
	masterTileList := []pixel.Rect{
		{pixel.V(1, 1), pixel.V(2, 2)},
		{pixel.V(3, 3), pixel.V(4, 4)},
		{pixel.V(5, 5), pixel.V(6, 6)},
		{pixel.V(7, 7), pixel.V(8, 8)},
		{pixel.V(9, 9), pixel.V(10, 10)},
		{pixel.V(11, 11), pixel.V(12, 12)},
	}

	sampleOccupiedList := []pixel.Rect{
		{pixel.V(1, 1), pixel.V(2, 2)},
		{pixel.V(3, 3), pixel.V(4, 4)},
	}

	expectedResult := []pixel.Rect{
		{pixel.V(5, 5), pixel.V(6, 6)},
		{pixel.V(7, 7), pixel.V(8, 8)},
		{pixel.V(9, 9), pixel.V(10, 10)},
		{pixel.V(11, 11), pixel.V(12, 12)},
	}

	actualResult := getAvailableTiles(masterTileList, sampleOccupiedList)

	match := true
	for i := range actualResult {
		if actualResult[i] != expectedResult[i] {
			match = false
		}
	}
	if !match {
		t.Error("Test failed: expected rect list does not match actual")
	}

}

func TestGetShotMoveDirection(t *testing.T) {
	if getShotMoveDirection(0)[0] != 0.25 || getShotMoveDirection(0)[1] != 0 {
		t.Error("Test failed: expected rect list does not match actual")
	}
}
