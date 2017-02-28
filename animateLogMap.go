package main

import (
	"image"
	"image/gif"
	"os"

	"github.com/plantimals/logisticmap/logisticmap"
)

type Animation struct {
	images []*image.Paletted
	delays []int
}

func main() {
	lm := logisticmap.NewLogisticMap()
	frames := 1
	imgs := make([]*image.Paletted, frames)
	delays := make([]int, frames)
	start := float64(3.65)
	stop := float64(3.75)
	for i := 0; i < frames; i++ {
		imgs[i] = lm.GetImage(start, stop, 0.000001)
		delays[i] = 5
		start += 0.01
		stop -= 0.01
	}
	output, _ := os.OpenFile("test2.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer output.Close()
	gif.EncodeAll(output, &gif.GIF{
		Image: imgs,
		Delay: delays,
	})
}
