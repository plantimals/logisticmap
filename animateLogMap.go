package main

import (
	"fmt"
	"os"

	"github.com/plantimals/logisticmap/logisticmap"
)

func main() {

	var xmin = 3.0
	var step = 0.005
	for i := 0; i < 200; i++ {
		OutputImage(i, xmin, xmin+step, 0.0, 1.0)
		xmin += step
	}
}

func OutputImage(index int, xmin float64, xmax float64, ymin float64, ymax float64) {
	output, _ := os.OpenFile(fmt.Sprintf("%v.png", index), os.O_WRONLY|os.O_CREATE, 0600)
	logisticmap.GetPNG(output, &logisticmap.Config{
		BurnIn:      100000,
		Take:        10000,
		Parallelism: 6,
		Scale:       50000,
		AspectRatio: 0.05,
		YMin:        ymin,
		YMax:        ymax,
		XMin:        xmin,
		XMax:        xmax,
	})
	output.Close()
}
