// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Run with "web" command-line argument for web server.
// See page 13.
//!+main

// Lissajous generates GIF animations of random Lissajous figures.
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
)

//!+main

var palette = []color.Color{color.Black, color.RGBA{0, 255, 0, 255}, color.RGBA{255, 0, 0, 255}}

const ( // first color in palette
	blackIndex = 0
	greenIndex = 1
	redIndex   = 2
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// Example query
// GET http://localhost:8000/?cycles=5&res=0.001&size=100&nframes=64&delay=8

func handler(w http.ResponseWriter, r *http.Request) {
	const (
		defaultCycles  = 5     // number of complete x oscillator revolutions
		defaultRes     = 0.001 // angular resolution
		defaultSize    = 100   // image canvas covers [-size..+size]
		defaultNframes = 64    // number of animation frames
		defaultDelay   = 8     // delay between frames in 10ms units
	)

	cycles := getFloatQueryParam(r, "cycles", defaultCycles)
	res := getFloatQueryParam(r, "res", defaultRes)
	size := getIntQueryParam(r, "size", defaultSize)
	nframes := getIntQueryParam(r, "nframes", defaultNframes)
	delay := getIntQueryParam(r, "delay", defaultDelay)

	lissajous(w, cycles, res, size, nframes, delay)
}

func getIntQueryParam(r *http.Request, param string, defaultValue int) int {
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Printf("Invalid value for %s: %s, using default: %d", param, valueStr, defaultValue)
		return defaultValue
	}
	return int(value)
}

func getFloatQueryParam(r *http.Request, param string, defaultValue float64) float64 {
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		log.Printf("Invalid value for %s: %s, using default: %f", param, valueStr, defaultValue)
		return defaultValue
	}
	return value
}

func lissajous(w io.Writer, cycles float64, res float64, size int, nframes int, delay int) {
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), greenIndex)
			img.SetColorIndex(size+int(x*float64(size)+0.7), size+int(y*float64(size)+0.7), redIndex)
		}

		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(w, &anim) // NOTE: ignoring encoding errors
}

//!-main
