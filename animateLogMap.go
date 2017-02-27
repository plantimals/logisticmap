package main

import (
	"image"
	"image/gif"
	"log"
	"os"
)

func NewLogisticMap() {
	x_dim = int(float64((XSTOP - XSTART)) / XSTEP)
	y_dim = int((x_dim * 3) / 4)
	parallelism := 8

	log.Printf("XSTEP: %v", XSTEP)
	log.Printf("XSTART: %v", XSTART)
	log.Printf("XSTOP: %v", XSTOP)
	log.Printf("x_dim: %v", x_dim)
	log.Printf("y_dim: %v", y_dim)

	slices := make(map[int]*VSlice)
	var fanout []<-chan *VSlice

	regions := paramGen(XSTART, XSTOP, XSTEP)

	for i := 0; i < parallelism; i++ {
		fanout = append(fanout, iterateGen(regions))
	}
	log.Printf("fanout size: %v", len(fanout))
	for vslice := range fanin(fanout) {
		slices[vslice.idx] = vslice
	}

	img := image.NewPaletted(image.Rect(0, 0, x_dim, y_dim), palette)
	var images []*image.Paletted
	fillImage(slices, img)
	images = append(images, img)
	output, err := os.OpenFile("test.gif", os.O_WRONLY|os.O_CREATE, 0600)
	handle(err)
	defer output.Close()
	gif.EncodeAll(output, &gif.GIF{
		Image: images,
		Delay: []int{1},
	})
}
