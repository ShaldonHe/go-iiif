package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

var (
	Input      string
	Output     string
	Number     int
	Alpha      int
	InputSize  int
	OutputSize int
	Mode       int
	V, VV      bool
)

func init() {
	flag.StringVar(&Input, "i", "", "input image path")
	flag.StringVar(&Output, "o", "", "output image path")
	flag.IntVar(&Number, "n", 0, "number of primitives")
	flag.IntVar(&Alpha, "a", 128, "alpha value")
	flag.IntVar(&InputSize, "r", 256, "resize large input images to this size")
	flag.IntVar(&OutputSize, "s", 1024, "output image size")
	flag.IntVar(&Mode, "m", 1, "0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect")
	flag.BoolVar(&V, "v", false, "verbose")
	flag.BoolVar(&VV, "vv", false, "very verbose")
}

func errorMessage(message string) bool {
	fmt.Fprintln(os.Stderr, message)
	return false
}

func main() {
	// parse and validate arguments
	flag.Parse()
	ok := true
	if Input == "" {
		ok = errorMessage("ERROR: input argument required")
	}
	if Output == "" {
		ok = errorMessage("ERROR: output argument required")
	}
	if Number == 0 {
		ok = errorMessage("ERROR: number argument required")
	}
	if !ok {
		fmt.Println("Usage: primitive [OPTIONS] -i input -o output -n shape_count")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// set log level
	if V {
		primitive.LogLevel = 1
	}
	if VV {
		primitive.LogLevel = 2
	}

	// seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// read input image
	primitive.Log(1, "reading %s\n", Input)
	input, err := primitive.LoadImage(Input)
	if err != nil {
		panic(err)
	}

	// scale down input image if needed
	size := uint(InputSize)
	input = resize.Thumbnail(size, size, input, resize.Bilinear)

	// determine output options
	ext := strings.ToLower(filepath.Ext(Output))
	saveFrames := strings.Contains(Output, "%") && ext != ".gif"

	// run algorithm
	model := primitive.NewModel(input, Alpha, OutputSize, primitive.Mode(Mode))
	start := time.Now()
	for i := 1; i <= Number; i++ {
		// find optimal shape and add it to the model
		model.Step()
		elapsed := time.Since(start).Seconds()
		primitive.Log(1, "iteration %d, time %.3f, score %.6f\n", i, elapsed, model.Score)

		// write output image(s)
		if saveFrames || i == Number {
			path := Output
			if saveFrames {
				path = fmt.Sprintf(Output, i)
			}
			primitive.Log(1, "writing %s\n", path)
			switch ext {
			case ".png":
				primitive.SavePNG(path, model.Context.Image())
			case ".svg":
				primitive.SaveFile(path, model.SVG())
			case ".gif":
				frames := model.Frames(0.001)
				primitive.SaveGIFImageMagick(path, frames, 50, 250)
			}
		}
	}
}
