package main

import (
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"
	"sync"
)

//var (
//	BURN_IN = 100000
//	TAKE    = 1000
//	XSTEP   = float64(0.0001)
//	XSTART  = float64(2.9)
//	XSTOP   = float64(4.0)
//	x_dim   int
//	y_dim   int
//)

func newVSlice(idx int, param float64) *VSlice {
	vSlice := new(VSlice)
	vSlice.idx = idx
	vSlice.param = param
	vSlice.levels = make([]float64, TAKE)
	return vSlice
}

var palette = []color.Color{
	color.RGBA{0x00, 0x00, 0x00, 0xff}, //black
	color.RGBA{0xff, 0xff, 0xff, 0xff}, //white
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type LogisicMap struct {
	burnIn      int
	take        int
	step        float64
	start       float64
	stop        float64
	x_dim       int
	y_dim       int
	parallelism int
}

const (
	burnIn      = 100000
	take        = 1000
	step        = 0.0001
	xStart      = float64(2.9)
	xStop       = float64(4.0)
	parallelism = 8
)

func NewLogisticMap() *LogisicMap {
	lm = new(LogisicMap)
	lm.burnIn = burnIn
	lm.take = take
	lm.start = start
	lm.stop = stop
	lm.step = step

	lm.x_dim = int(float64((XSTOP - XSTART)) / XSTEP)
	lm.y_dim = int((x_dim * 3) / 4)
	lm.parallelism = parallelism
	return lm
}

func (lm *LogisicMap) Parallelism(p int) {
	lm.parallelism = p
}

func (lm *LogisicMap) GetFrame(start float64, stop float64, step float64) {

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

func fillImage(slices map[int]*VSlice, img *image.Paletted) {
	yf := float64(y_dim)
	for x := int(0); x < x_dim; x++ {
		if _, ok := slices[x]; ok {
			for _, p := range slices[x].levels {
				y := int((1 - p) * yf)
				img.Set(x, y, palette[1])
			}
		}
	}
}

type VSlice struct {
	idx    int
	param  float64
	levels []float64
}

func fanin(fanout []<-chan *VSlice) <-chan *VSlice {
	var wg sync.WaitGroup
	out := make(chan *VSlice)
	output := func(c <-chan *VSlice) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(fanout))
	for _, c := range fanout {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func paramGen(start float64, stop float64, step float64) <-chan *VSlice {
	out := make(chan *VSlice)
	go func() {
		sliceCount := int((stop - start) / step)
		log.Printf("sliceCount: %v", sliceCount)
		for i := 0; i < sliceCount; i++ {
			p := start + (float64(i) * step)
			out <- newVSlice(i, p)
		}
		close(out)
	}()
	return out
}

func iterateGen(in <-chan *VSlice) <-chan *VSlice {
	out := make(chan *VSlice)
	go func() {
		for vslice := range in {
			iterate(vslice)
			out <- vslice
		}
		close(out)
	}()
	return out
}

func iterate(vslice *VSlice) {
	var x = float64(0.7)
	param := vslice.param
	for i := 0; i < BURN_IN; i++ {
		x = (param * x) * (1 - x)
	}
	for i := 0; i < TAKE; i++ {
		x = (param * x) * (1 - x)
		vslice.levels[i] = x
	}
}
