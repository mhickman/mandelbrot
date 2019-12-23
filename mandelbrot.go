package mandelbrot

import (
	"image"
	"image/color"
	"math"
	"sync"
)

type pointLocation complex128

// If |current|^2 > max then we know that the point is _not_ in
// the Mandelbrot set.
const max float64 = 4.0
const maxIterations = 5000

func newPointLocation(r float64, i float64) pointLocation {
	return pointLocation(complex(r, i))
}

type Point struct {
	location  pointLocation
	iteration int64
	current   pointLocation

	inSet     bool
	processed bool
}

func (c pointLocation) abs2() float64 {
	r := real(complex128(c))
	i := imag(complex128(c))
	return r*r + i*i
}

// Take the point through one iteration of Z^2 + C
func (p *Point) Iterate() {
	p.iteration = p.iteration + 1
	p.current = p.current*p.current + p.location
}

func (p *Point) IsMandelbrot() bool {
	return p.inSet
}

func (p *Point) DetermineMembership() bool {
	if p.processed {
		return p.inSet
	}

	i := 0

	for i < maxIterations && p.current.abs2() < max {
		p.Iterate()
		i++
	}

	if p.current.abs2() < max {
		p.inSet = true
	} else {
		p.inSet = false
	}

	p.processed = true

	return p.inSet
}

// Returns a new point at (r, i) with 0 iterations done so far.
func NewPoint(r float64, i float64) Point {
	return Point{
		location:  newPointLocation(r, i),
		iteration: 0,
		current:   newPointLocation(0, 0),
		processed: false,
	}
}

type Grid struct {
	center        complex128
	width, height int64
	pixelWidth    float64

	// [0][0] is bottom left
	// [width-1]][0] is bottom right
	// [0][height-1] is top left
	// [width-1][height-1] is top right
	points [][]*Point
}

func NewGrid(center complex128, width, height int64, pixelWidth float64) Grid {
	points := make([][]*Point, width)

	halfWidth := 0.5 * pixelWidth * float64(width)
	halfHeight := 0.5 * pixelWidth * float64(height)

	bottomLeft := center - complex(halfWidth, halfHeight)

	for col := range points {
		points[col] = make([]*Point, height)

		for row := range points[col] {
			point := NewPoint(
				real(bottomLeft)+float64(col)*pixelWidth,
				imag(bottomLeft)+float64(row)*pixelWidth,
			)

			points[col][row] = &point
		}
	}

	return Grid{
		center:     center,
		width:      width,
		height:     height,
		pixelWidth: pixelWidth,
		points:     points,
	}
}

func (g *Grid) Points() [][]*Point {
	return g.points
}

func (g *Grid) IterateAll() {
	var wg sync.WaitGroup

	for _, row := range g.points {
		for _, point := range row {
			wg.Add(1)

			go func(point *Point) {
				defer wg.Done()
				point.DetermineMembership()
			}(point)
		}
	}

	wg.Wait()
}

func (g *Grid) GenerateImage() image.Image {
	r := image.Rect(0, 0, len(g.points), len(g.points[0]))

	im := image.NewCMYK(r)

	min := int64(math.MaxInt64)
	max := int64(0)

	for _, row := range g.points {
		for _, point := range row {
			if !point.inSet {
				if min > point.iteration {
					min = point.iteration
				}

				if max < point.iteration {
					max = point.iteration
				}
			}
		}
	}

	for i, row := range g.points {
		for j, point := range row {
			if point.inSet {
				im.Set(i, j, color.Black)
			} else {
				green := color.RGBA{
					R: 0,
					G: 0xff,
					B: 0,
					A: 0xff,
				}

				red := color.RGBA{
					R: 0xff,
					G: 0,
					B: 0,
					A: 0xff,
				}

				//grey := color.RGBA{
				//	R: 0x30,
				//	G: 0x30,
				//	B: 0x30,
				//	A: 0xff,
				//}

				im.Set(i, j, interpolateColors(green, red, float64(point.iteration) / float64(max)))
			}
		}
	}

	return im
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

	return color.RGBA {
		R: interpolateInt(r1, r2, a),
		G: interpolateInt(g1, g2, a),
		B: interpolateInt(b1, b2, a),
		A: interpolateInt(a1, a2, a),
	}
}
