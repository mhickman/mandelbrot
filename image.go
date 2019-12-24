package mandelbrot

import (
	"image/color"
	"math"
)

type ColorPalette interface {
	Color(point Point) color.Color
}

type linearColorPalette struct {
	lowColor      color.Color
	highColor     color.Color
	inSetColor    color.Color
	maxIterations int64
}

func NewLinearPalette(
	grid *Grid,
	lowColor color.Color,
	highColor color.Color,
	inSetColor color.Color) ColorPalette {
	max := int64(0)

	for _, row := range grid.points {
		for _, point := range row {
			if !point.inSet {
				if max < point.iteration {
					max = point.iteration
				}
			}
		}
	}

	return &linearColorPalette{
		lowColor:      lowColor,
		highColor:     highColor,
		inSetColor:    inSetColor,
		maxIterations: max,
	}
}

func interpolateInt(a uint8, b uint8, p float64) uint8 {
	p = math.Sqrt(p)
	compP := 1.0 - p

	aFloat := p * float64(a)
	bFloat := compP * float64(b)

	return uint8(aFloat + bFloat)
}

func interpolateColors(color1 color.Color, color2 color.Color, a float64) color.Color {
	color1r, color1g, color1b, color1a := color1.RGBA()
	color2r, color2g, color2b, color2a := color2.RGBA()

	var r1, r2, g1, g2, b1, b2, a1, a2 uint8

	r1 |= uint8(color1r >> 8)
	r2 |= uint8(color2r >> 8)
	g1 |= uint8(color1g >> 8)
	g2 |= uint8(color2g >> 8)
	b1 |= uint8(color1b >> 8)
	b2 |= uint8(color2b >> 8)
	a1 |= uint8(color1a >> 8)
	a2 |= uint8(color2a >> 8)

	return color.RGBA{
		R: interpolateInt(r1, r2, a),
		G: interpolateInt(g1, g2, a),
		B: interpolateInt(b1, b2, a),
		A: interpolateInt(a1, a2, a),
	}
}

func (cp *linearColorPalette) Color(point Point) color.Color {
	if point.inSet {
		return cp.inSetColor
	} else {
		return interpolateColors(
			cp.lowColor,
			cp.highColor,
			float64(point.iteration)/float64(cp.maxIterations),
		)
	}
}
