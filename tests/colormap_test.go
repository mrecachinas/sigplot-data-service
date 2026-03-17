package main

import (
	"testing"

	"github.com/spectriclabs/sigplot-data-service/internal/image"
)

func TestGetColorControlPoints(t *testing.T) {
	expected := []struct {
		Input    string
		NumElems int
	}{
		{Input: "Greyscale", NumElems: 3},
		{Input: "Ramp Colormap", NumElems: 7},
		{Input: "Color Wheel", NumElems: 7},
		{Input: "Spectrum", NumElems: 7},
		{Input: "calewhite", NumElems: 7},
		{Input: "HotDesat", NumElems: 8},
		{Input: "Sunset", NumElems: 7},
	}
	for _, exp := range expected {
		result := image.GetColorControlPoints(exp.Input)
		if len(result) != exp.NumElems {
			t.Errorf(
				"GetColorControlPoints(%s) returned %d elements, expected %d",
				exp.Input,
				len(result),
				exp.NumElems,
			)
		}
	}

	// Unknown colormap should fall back to default (Ramp Colormap = 7 points)
	result := image.GetColorControlPoints("nonexistent")
	if len(result) != 7 {
		t.Errorf("GetColorControlPoints(nonexistent) returned %d elements, expected 7 (default)", len(result))
	}
}

func TestMakeColorPalette(t *testing.T) {
	controlColors := image.GetColorControlPoints("Greyscale")
	numColors := 10
	result := image.MakeColorPalette(controlColors, numColors)
	if len(result) != numColors {
		t.Errorf("MakeColorPalette returned %d colors, expected %d", len(result), numColors)
	}
}
