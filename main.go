// This program sends GIFs that display the movement
// of a lissajous function as the HTTP response
// body of a server.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var palette = []color.Color{color.Black, color.RGBA{0x00, 0xff, 0x00, 1}, color.RGBA{0xff, 0x00, 0x00, 1}, color.RGBA{0x00, 0x00, 0xff, 1}}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

const (
	defaultCycles  = 5     // number of compete x oscillator revolutions
	defaultRes     = 0.001 // angular resolution
	defaultSize    = 100   // image canvas covers [-size..+size]
	defaultNframes = 64    // number of animation frames
	defaultDelay   = 8     // delay between frames in 10ms units
)

// handler writes the lissajous GIF to the response.
func handler(w http.ResponseWriter, r *http.Request) {
	var cycles, size, nframes, delay int = defaultCycles, defaultSize, defaultNframes, defaultDelay
	var res float64 = defaultRes
	var err error
	query := r.URL.Query()
	if _, ok := query["cycles"]; ok {
		if cycles, err = strconv.Atoi(query.Get("cycles")); err != nil {
			cycles = defaultCycles
		}
	}
	if _, ok := query["size"]; ok {
		if size, err = strconv.Atoi(query.Get("size")); err != nil {
			size = defaultSize
		}
	}
	if _, ok := query["nframes"]; ok {
		if nframes, err = strconv.Atoi(query.Get("nframes")); err != nil {
			nframes = defaultNframes
		}
	}
	if _, ok := query["delay"]; ok {
		if delay, err = strconv.Atoi(query.Get("delay")); err != nil {
			delay = defaultDelay
		}
	}
	if _, ok := query["res"]; ok {
		if res, err = strconv.ParseFloat(query.Get("res"), 64); err != nil {
			res = defaultRes
		}
	}
	lissajous(w, cycles, size, nframes, delay, res)
}

func lissajous(out io.Writer, cycles, size, nframes, delay int, res float64) {
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), uint8(rand.Intn(3)+1))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
