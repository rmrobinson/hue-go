package hue

import (
	"testing"
)

type colorTest struct {
	rgb      RGB
	xyGamutA XY
	xyGamutB XY
	xyGamutC XY
	ct       uint16
}

// Sample mappings retrieved from http://www.developers.meethue.com/documentation/hue-xy-values

var colorAliceBlue = colorTest{
	rgb: RGB{
		Red:   239,
		Green: 247,
		Blue:  255,
	},
	xyGamutA: XY{X: 0.3088, Y: 0.3212},
	xyGamutB: XY{X: 0.3092, Y: 0.321},
	xyGamutC: XY{X: 0.3088, Y: 0.3212},
}

var colorChartreuse = colorTest{
	rgb: RGB{
		Red:   127,
		Green: 255,
		Blue:  0,
	},
	xyGamutA: XY{X: 0.2682, Y: 0.6632},
	xyGamutB: XY{X: 0.408, Y: 0.517},
	xyGamutC: XY{X: 0.2505, Y: 0.6395},
}

func TestXY_FromRGB(t *testing.T) {
	var output XY

	output.FromRGB(colorAliceBlue.rgb, "LST001")

	if output.X != colorAliceBlue.xyGamutA.X || output.Y != colorAliceBlue.xyGamutA.Y {
		t.Errorf("Incorrect conversion of alice blue from RGB to gamut A, expected [%f,%f], got [%f,%f]\n", colorAliceBlue.xyGamutA.X, colorAliceBlue.xyGamutA.Y, output.X, output.Y)
	}

	output.FromRGB(colorAliceBlue.rgb, "LCT001")

	if output.X != colorAliceBlue.xyGamutB.X || output.Y != colorAliceBlue.xyGamutB.Y {
		t.Errorf("Incorrect conversion of alice blue from RGB to gamut B, expected [%f,%f], got [%f,%f]\n", colorAliceBlue.xyGamutB.X, colorAliceBlue.xyGamutB.Y, output.X, output.Y)
	}

	output.FromRGB(colorAliceBlue.rgb, "LLC020")

	if output.X != colorAliceBlue.xyGamutC.X || output.Y != colorAliceBlue.xyGamutC.Y {
		t.Errorf("Incorrect conversion of alice blue from RGB to gamut C, expected [%f,%f], got [%f,%f]\n", colorAliceBlue.xyGamutC.X, colorAliceBlue.xyGamutC.Y, output.X, output.Y)
	}

	output.FromRGB(colorChartreuse.rgb, "LST001")

	if output.X != colorChartreuse.xyGamutA.X || output.Y != colorChartreuse.xyGamutA.Y {
		t.Errorf("Incorrect conversion of chartreuse from RGB to gamut A, expected [%f,%f], got [%f,%f]\n", colorChartreuse.xyGamutA.X, colorChartreuse.xyGamutA.Y, output.X, output.Y)
	}

	output.FromRGB(colorChartreuse.rgb, "LCT001")

	if output.X != colorChartreuse.xyGamutB.X || output.Y != colorChartreuse.xyGamutB.Y {
		t.Errorf("Incorrect conversion of chartreuse from RGB to gamut B, expected [%f,%f], got [%f,%f]\n", colorChartreuse.xyGamutB.X, colorChartreuse.xyGamutB.Y, output.X, output.Y)
	}

	output.FromRGB(colorChartreuse.rgb, "LLC020")

	if output.X != colorChartreuse.xyGamutC.X || output.Y != colorChartreuse.xyGamutC.Y {
		t.Errorf("Incorrect conversion of chartreuse from RGB to gamut C, expected [%f,%f], got [%f,%f]\n", colorChartreuse.xyGamutC.X, colorChartreuse.xyGamutC.Y, output.X, output.Y)
	}
}

func TestRGB_FromXY(t *testing.T) {
	var output RGB

	output.FromXY(colorAliceBlue.xyGamutA, "LST001")

	if output.Red != colorAliceBlue.rgb.Red || output.Green != colorAliceBlue.rgb.Green || output.Blue != colorAliceBlue.rgb.Blue {
		t.Errorf("Incorrect conversion of alice blue from gamut A to RGB, expected [%d,%d,%d], got [%d,%d,%d]\n", colorAliceBlue.rgb.Red, colorAliceBlue.rgb.Green, colorAliceBlue.rgb.Blue, output.Red, output.Green, output.Blue)
	}

	output.FromXY(colorAliceBlue.xyGamutB, "LCT001")

	if output.Red != colorAliceBlue.rgb.Red || output.Green != colorAliceBlue.rgb.Green || output.Blue != colorAliceBlue.rgb.Blue {
		t.Errorf("Incorrect conversion of alice blue from gamut B to RGB, expected [%d,%d,%d], got [%d,%d,%d]\n", colorAliceBlue.rgb.Red, colorAliceBlue.rgb.Green, colorAliceBlue.rgb.Blue, output.Red, output.Green, output.Blue)
	}

	output.FromXY(colorAliceBlue.xyGamutC, "LLC020")

	if output.Red != colorAliceBlue.rgb.Red || output.Green != colorAliceBlue.rgb.Green || output.Blue != colorAliceBlue.rgb.Blue {
		t.Errorf("Incorrect conversion of alice blue from gamut C to RGB, expected [%d,%d,%d], got [%d,%d,%d]\n", colorAliceBlue.rgb.Red, colorAliceBlue.rgb.Green, colorAliceBlue.rgb.Blue, output.Red, output.Green, output.Blue)
	}

	output.FromXY(colorChartreuse.xyGamutA, "LST001")

	if output.Red != colorChartreuse.rgb.Red || output.Green != colorChartreuse.rgb.Green || output.Blue != colorChartreuse.rgb.Blue {
		t.Errorf("Incorrect conversion of chartreuse from gamut A to RGB, expected [%d,%d,%d], got [%d,%d,%d]\n", colorChartreuse.rgb.Red, colorChartreuse.rgb.Green, colorChartreuse.rgb.Blue, output.Red, output.Green, output.Blue)
	}

	output.FromXY(colorChartreuse.xyGamutB, "LCT001")

	if output.Red != colorChartreuse.rgb.Red || output.Green != colorChartreuse.rgb.Green || output.Blue != colorChartreuse.rgb.Blue {
		t.Errorf("Incorrect conversion of chartreuse from gamut B to RGB, expected [%d,%d,%d], got [%d,%d,%d]\n", colorChartreuse.rgb.Red, colorChartreuse.rgb.Green, colorChartreuse.rgb.Blue, output.Red, output.Green, output.Blue)
	}

	output.FromXY(colorChartreuse.xyGamutC, "LLC020")

	if output.Red != colorChartreuse.rgb.Red || output.Green != colorChartreuse.rgb.Green || output.Blue != colorChartreuse.rgb.Blue {
		t.Errorf("Incorrect conversion of chartreuse from gamut C to RGB, expected [%d,%d,%d], got [%d,%d,%d]\n", colorChartreuse.rgb.Red, colorChartreuse.rgb.Green, colorChartreuse.rgb.Blue, output.Red, output.Green, output.Blue)
	}

}
