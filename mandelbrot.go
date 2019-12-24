package mandelbrot

import (
	"image"
	"image/color"
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

func (c pointLocation) abs2() float64 {
	r := real(complex128(c))
	i := imag(complex128(c))
	return r*r + i*i
}

type Point struct {
	location  pointLocation
	iteration int64
	current   pointLocation

	inSet     bool
	processed bool
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

func (g *Grid) GenerateImageWithPalette(p ColorPalette) image.Image {
	r := image.Rect(0, 0, len(g.points), len(g.points[0]))

	im := image.NewCMYK(r)

	for i, row := range g.points {
		for j, point := range row {
			im.Set(i, j, p.Color(*point))
		}
	}

	return im
}

func (g *Grid) GenerateImage() image.Image {
	green := color.RGBA{
		R: 0,
		G: 0xff,
		B: 0,
		A: 0xff,
	}
	p := NewLinearPalette(g, color.Black, green, color.Black)
	return g.GenerateImageWithPalette(p)
}
