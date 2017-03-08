package gui

import (
	"math"

	"github.com/mattn/go-gtk/gdk"
)

func titleBlue() *gdk.Color {
	return colorToGDKColor(40, 128, 185)
}

func errorRed() *gdk.Color {
	return colorToGDKColor(255, 180, 180)
}

func colorToGDKColor(red, green, blue uint8) *gdk.Color {
	return gdk.NewColorRGB(
		uint16(math.Pow(float64(red), 2)),
		uint16(math.Pow(float64(green), 2)),
		uint16(math.Pow(float64(blue), 2)),
	)
}
