package mandelbrot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPoint(t *testing.T) {
	newPoint := NewPoint(1.0, 4.0)

	assert.Equal(t, newPointLocation(0.0, 0.0), newPoint.current)
	assert.Equal(t, newPointLocation(1.0, 4.0), newPoint.location)
	assert.Equal(t, int64(0), newPoint.iteration)
}

func Test_abs2(t *testing.T) {
	newPoint := newPointLocation(2.0, 3.0)
	assert.Equal(t, 13.0, newPoint.abs2())
}

func TestPoint_Iterate(t *testing.T) {
	newPoint := NewPoint(1.0, 2.0)

	newPoint.Iterate()
	assert.Equal(t, int64(1), newPoint.iteration)
	assert.Equal(t, newPointLocation(1.0, 2.0), newPoint.current)

	newPoint.Iterate()
	assert.Equal(t, int64(2), newPoint.iteration)
	assert.Equal(t, newPointLocation(-2.0, 6.0), newPoint.current)
}

func TestPoint_DetermineMembership(t *testing.T) {
	inSetPoint := NewPoint(.23, 0)
	inSet := inSetPoint.DetermineMembership()

	assert.True(t, inSet)
	assert.True(t, inSetPoint.inSet)

	outSetPoint := NewPoint(.26, 0)
	inSet = outSetPoint.DetermineMembership()

	assert.False(t, inSet)
	assert.False(t, outSetPoint.inSet)
}

func TestNewGrid(t *testing.T) {
	var fourGrid Grid

	assert.NotPanics(t, func() {
		fourGrid = NewGrid(complex(0.0, 0.0), 2, 2, 1.0)
	})

	assert.Equal(t, 2, len(fourGrid.points))

	for _, col := range fourGrid.points {
		assert.Equal(t, 2, len(col))
	}
}

func TestGrid_IterateAll(t *testing.T) {
	grid := NewGrid(complex(0.0, 0.0), 2, 2, 1.0)

	grid.IterateAll()

	for _, row := range grid.points {
		for _, point := range row {
			assert.True(t, point.processed)
		}
	}
}
