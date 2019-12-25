package mandelbrot

import (
	"github.com/stretchr/testify/assert"
	"image/color"
	"testing"
)

var red, green, blue color.RGBA

func init() {
	red = color.RGBA{
		R: 0xff,
		G: 0,
		B: 0,
		A: 0xff,
	}

	green = color.RGBA{
		R: 0,
		G: 0xff,
		B: 0,
		A: 0,
	}

	blue = color.RGBA{
		R: 0,
		G: 0,
		B: 0xff,
		A: 0,
	}
}

func TestNewMultiColorGradient_sortedColors(t *testing.T) {
	lowestColor := GradientColor{
		Percent: 0.25,
		Color:   color.Black,
	}

	highestColor := GradientColor{
		Percent: 0.75,
		Color:   color.White,
	}

	unsortedColors := []GradientColor{highestColor, lowestColor}
	grid := NewGrid(complex(0.0, 0.0), 2, 2, 1.0)

	gradient := NewMultiColorGradient(
		&grid,
		unsortedColors,
		red,
		green,
		blue,
	)

	// Force casting to underlying type to make assertions about
	// the underlying struct.
	multiGradient, ok := gradient.(*multiColorGradient)

	if ok {
		assert.Equal(t, red, multiGradient.colors[0].Color)
		assert.Equal(t, color.Black, multiGradient.colors[1].Color)
		assert.Equal(t, color.White, multiGradient.colors[2].Color)
		assert.Equal(t, green, multiGradient.colors[3].Color)
	} else {
		t.Fail()
	}
}
