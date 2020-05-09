package logisticmap

import (
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"io"
	"log"
	"sync"
)

type Config struct {
	BurnIn      int
	Take        int
	Parallelism int
	Scale       int
	AspectRatio float64
	YMin        float64
	YMax        float64
	yRange      float64
	XMin        float64
	XMax        float64
	step        float64
	pX          int
	pY          int
}

func newVSlice(idx int, param float64, yMin float64, yMax float64, take int) *VSlice {
	vSlice := new(VSlice)
	vSlice.idx = idx
	vSlice.param = param
	vSlice.levels = make([]float64, take)
	vSlice.yMin = yMin
	vSlice.yMax = yMax
	return vSlice
}

type VSlice struct {
	idx    int
	param  float64
	levels []float64
	lvlIdx int
	yMin   float64
	yMax   float64
}

func (vs *VSlice) Fill() bool {
	return vs.lvlIdx < len(vs.levels)
}

func (vs *VSlice) Add(y float64) bool {
	answer := false
	if y > vs.yMin && y < vs.yMax {
		vs.levels[vs.lvlIdx] = y
		vs.lvlIdx++
		answer = true
	}
	return answer
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

func GetGIF(writer io.Writer, config *Config) {
	config.yRange = config.YMax - config.YMin
	config.pX = int(config.AspectRatio * float64(config.Scale))
	config.pY = int(config.Scale)
	config.step = (config.XMax - config.XMin) / float64(config.pX)

	img := GetImage(config)
	var images []*image.Paletted
	images = append(images, img)
	gif.EncodeAll(writer, &gif.GIF{
		Image: images,
	})
	log.Printf("completed writing image")
}

func GetPNG(writer io.Writer, config *Config) {
	config.yRange = config.YMax - config.YMin
	config.pX = int(config.AspectRatio * float64(config.Scale))
	config.pY = int(config.Scale)
	config.step = (config.XMax - config.XMin) / float64(config.pX)

	img := GetImage(config)
	png.Encode(writer, img)
	log.Printf("completed writing image")
}

func Pan(writer io.Writer, config *Config, dx float64, dy float64, frames int, delay int) {
	config.yRange = config.YMax - config.YMin
	config.pX = int(config.AspectRatio * float64(config.Scale))
	config.pY = int(config.Scale)
	config.step = (config.XMax - config.XMin) / float64(config.pX)
	var images []*image.Paletted
	var delays []int
	for idx := 0; idx < frames; idx++ {
		config.XMin += dx
		config.XMax += dx
		config.YMin += dy
		config.YMax += dy
		images = append(images, GetImage(config))
		delays = append(delays, delay)
		log.Printf("completed image %v", idx)
	}
	log.Printf("len(images): %v", len(images))
	gif.EncodeAll(writer, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

func GetImage(config *Config) *image.Paletted {

	slices := make(map[int]*VSlice)
	var fanout []<-chan *VSlice

	regions := paramGen(config)

	for i := 0; i < config.Parallelism; i++ {
		fanout = append(fanout, iterateGen(regions, config.BurnIn, config.Take))
	}
	for vslice := range fanin(fanout) {
		slices[vslice.idx] = vslice
	}

	img := image.NewPaletted(image.Rect(0, 0, config.pX, config.pY), palette)
	fillImage(slices, img, config)
	return img
}

func fillImage(slices map[int]*VSlice, img *image.Paletted, config *Config) {
	log.Printf("%v slices\n", len(slices))
	log.Printf("%v points per slice\n", len(slices[0].levels))
	yf := float64(config.pY)
	for x := int(0); x < config.pX; x++ {
		if _, ok := slices[x]; ok {
			for _, p := range slices[x].levels {
				// this line flips the image, shifts it to fit the output, then scales it up
				y := int((((1 - p) - (1 - config.YMax)) / config.yRange) * yf)
				img.Set(x, int(y), palette[1])
			}
		}
	}
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

func paramGen(config *Config) <-chan *VSlice {
	start := config.XMin
	stop := config.XMax
	step := config.step
	out := make(chan *VSlice)
	go func() {
		sliceCount := int((stop - start) / step)
		for i := 0; i < sliceCount; i++ {
			p := start + (float64(i) * step)
			out <- newVSlice(i, p, config.YMin, config.YMax, config.Take)
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
	count := 0
	for vslice.Fill() {
		x = (param * x) * (1 - x)
		vslice.Add(x)
		if count > take*10 {
			break
		}
		count++
	}
}
