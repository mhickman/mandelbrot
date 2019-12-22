package main

import (
	"fmt"
	"github.com/mhickman/mandelbrot"
	"image/png"
	"os"
)

func main() {
	grid := mandelbrot.NewGrid(complex(0.0, 0.0), 2000, 1000, 0.002)
	grid.IterateAll()

	img := grid.GenerateImage()

	out, err := os.Create("./output.png")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Generated image to output.png")
}
