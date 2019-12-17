package mandelbrot

type pointLocation complex128

// If |current|^2 > max then we know that the point is _not_ in
// the Mandelbrot set.
const max float64 = 4.0
const maxIterations = 10_000

func newPointLocation(r float64, i float64) pointLocation {
	return pointLocation(complex(r, i))
}

type Point struct {
	location  pointLocation
	iteration int64
	current   pointLocation

	inSet bool
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

func (p *Point) DetermineMembership() bool {
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

	return p.inSet
}

// Returns a new point at (r, i) with 0 iterations done so far.
func NewPoint(r float64, i float64) Point {
	return Point{
		location:  newPointLocation(r, i),
		iteration: 0,
		current:   newPointLocation(0, 0),
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
	points [][]Point
}

func NewGrid(center complex128, width, height int64, pixelWidth float64) Grid {
	points := make([][]Point, width)

	halfWidth := 0.5 * pixelWidth * float64(width)
	halfHeight := 0.5 * pixelWidth * float64(height)

	bottomLeft := center - complex(halfWidth, halfHeight)

	for col := range points {
		points[col] = make([]Point, height)

		for row := range points[col] {
			points[col][row] = NewPoint(
				real(bottomLeft)+float64(col)*pixelWidth,
				imag(bottomLeft)+float64(row)*pixelWidth,
			)
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
