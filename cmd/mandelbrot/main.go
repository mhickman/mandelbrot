package main

import (
	"fmt"
	"github.com/mhickman/mandelbrot"
	"image/png"
	"os"
)

func main() {
	grid := mandelbrot.NewGrid(complex(-0.743643887037151,  0.131825904205330), 5000, 5000, 0.000000001)
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
