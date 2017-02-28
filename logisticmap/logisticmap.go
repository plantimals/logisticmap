package logisticmap

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"sync"
)

func newVSlice(idx int, param float64, take int) *VSlice {
	vSlice := new(VSlice)
	vSlice.idx = idx
	vSlice.param = param
	vSlice.levels = make([]float64, take)
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
	xRes        int
	yRes        int
}

const (
	burnIn      = 100000
	take        = 1000
	parallelism = 8
)

func NewLogisticMap(xRes int, yRes int) *LogisicMap {
	lm := new(LogisicMap)
	lm.burnIn = burnIn
	lm.take = take
	lm.xRes = xRes
	lm.yRes = yRes

	lm.parallelism = parallelism
	return lm
}

func (lm *LogisicMap) Parallelism(p int) {
	lm.parallelism = p
}

func (lm *LogisicMap) GetGIF(writer io.Writer, start float64, stop float64, step float64) {
	img := lm.GetImage(start, stop, step)
	var images []*image.Paletted
	images = append(images, img)
	lm.WriteGIF(writer, images, []int{1})
}

func (lm *LogisicMap) WriteGIF(writer io.Writer, images []*image.Paletted, delays []int) {
	gif.EncodeAll(writer, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

type Config struct {
	Scale       int
	AspectRatio float64
	yMin        float64
	yMax        float64
	xMin        float64
	xMax        float64
}

func (lm *LogisicMap) Pan(writer io.Writer, config *Config, dx float64, dy float64, frames int) {
	pX := config.AspectRatio * config.Scale
	pY := scale

}

func (lm *LogisicMap) GetImage(start float64, stop float64, step float64) *image.Paletted {
	lm.start = start
	lm.stop = stop
	lm.step = step
	lm.x_dim = int(float64((stop - start)) / step)
	lm.y_dim = int((lm.x_dim * 3) / 4)
	slices := make(map[int]*VSlice)
	var fanout []<-chan *VSlice

	regions := paramGen(start, stop, step, lm.take)

	for i := 0; i < lm.parallelism; i++ {
		fanout = append(fanout, iterateGen(regions, burnIn, take))
	}
	for vslice := range fanin(fanout) {
		slices[vslice.idx] = vslice
	}

	img := image.NewPaletted(image.Rect(0, 0, lm.x_dim, lm.y_dim), palette)
	lm.fillImage(slices, img)
	return img
}

func (lm *LogisicMap) fillImage(slices map[int]*VSlice, img *image.Paletted) {
	yf := float64(lm.y_dim)
	for x := int(0); x < lm.x_dim; x++ {
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

func paramGen(start float64, stop float64, step float64, take int) <-chan *VSlice {
	out := make(chan *VSlice)
	go func() {
		sliceCount := int((stop - start) / step)
		for i := 0; i < sliceCount; i++ {
			p := start + (float64(i) * step)
			out <- newVSlice(i, p, take)
		}
		close(out)
	}()
	return out
}

func iterateGen(in <-chan *VSlice, burnIn int, take int) <-chan *VSlice {
	out := make(chan *VSlice)
	go func() {
		for vslice := range in {
			iterate(vslice, burnIn, take)
			out <- vslice
		}
		close(out)
	}()
	return out
}

func iterate(vslice *VSlice, burnIn int, take int) {
	var x = float64(0.7)
	param := vslice.param
	for i := 0; i < burnIn; i++ {
		x = (param * x) * (1 - x)
	}
	for i := 0; i < take; i++ {
		x = (param * x) * (1 - x)
		vslice.levels[i] = x
	}
}
