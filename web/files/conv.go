package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/disintegration/imaging"
)

func main() {

	src, err := imaging.Open("jeffrey-hamilton-7LYmnzX1buU-unsplash.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	// Resize srcImage to size = 128x128px using the Lanczos filter.
	g := src.Bounds()

	// Get height and width
	height := float64(g.Dy())
	width := float64(g.Dx())
	perc := float64(4999999) / float64(5480315)
	fmt.Printf(" %f %f \n perc= %f\n", width, height, perc)
	height = height * perc
	width = width * perc
	fmt.Printf(" %f %f\n", width, height)
	h := int(height)
	w := int(width)
	fmt.Printf(" %d %d\n", w, h)

	arg1 := "jeffrey-hamilton-7LYmnzX1buU-unsplash.jpg"
	arg2 := "jeffrey-hamilton-7LYmnzX1buU-unsplash.jpg"
	d := fmt.Sprintf("ffmpeg -y -i %s -vf scale=%d:%d %s", arg1, w, h, arg2)
	lsCmd := exec.Command("bash", "-c", d)
	_, err = lsCmd.Output()
	if err != nil {
		panic(err)
	}

}
