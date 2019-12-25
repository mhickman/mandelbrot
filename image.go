package mandelbrot

import (
	"image/color"
	"math"
	"sort"
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

type GradientColor struct {
	// A number in [0.0, 1.0] that describes at what point this color will
	// be 100%.
	Percent float64

	Color color.Color
}

type multiColorGradient struct {
	colors     []GradientColor
	inSetColor color.Color
}

func (g *multiColorGradient) Color(point Point) color.Color {
	panic("Implement me!")
}

func NewMultiColorGradient(
	colors []GradientColor,
	minColor color.Color,
	maxColor color.Color,
	inSetColor color.Color,
) ColorPalette {
	colors = append(colors, GradientColor{
		Percent: -0.01,
		Color:   minColor,
	})

	colors = append(colors, GradientColor{
		Percent: 1.01,
		Color:   maxColor,
	})

	sort.Slice(colors, func(i, j int) bool {
		return colors[i].Percent < colors[j].Percent
	})

	return &multiColorGradient{
		colors:     colors,
		inSetColor: inSetColor,
	}
}
