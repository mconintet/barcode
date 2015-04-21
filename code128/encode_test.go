package code128

import (
	"image/png"
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	// 34
	// HI3456
	// HI345678
	// HI345678HI
	// HI3456HI
	// HI3456789HI
}

func TestNewPng(t *testing.T) {
	ts := "HI3456HI"

	img, err := Encode([]byte(ts), 0, 0, 0, 2)

	if err != nil {
		t.Fatal(err)
	}

	if file, err := os.OpenFile(ts+".png", os.O_RDWR|os.O_CREATE, 0666); err != nil {
		t.Fatal(err)
	} else {
		png.Encode(file, img)
	}
}
