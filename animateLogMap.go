package main

import (
	"os"

	"github.com/plantimals/logisticmap/logisticmap"
)

func main() {
	output, _ := os.OpenFile("pan.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer output.Close()
	logisticmap.Pan(output, &logisticmap.Config{
		BurnIn:      100000,
		Take:        1000,
		Parallelism: 8,
		Scale:       1200,
		AspectRatio: 1.0,
		YMin:        0.0,
		YMax:        1.0,
		XMin:        3.65,
		XMax:        3.7,
	}, 0.00006, 0.0, 1000, 5)

}
