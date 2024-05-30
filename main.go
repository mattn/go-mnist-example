package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/mattn/go-tflite"
	"github.com/nfnt/resize"
)

func top(a []float32) int {
	t := 0
	m := float32(0)
	for i, e := range a {
		if i == 0 || e > m {
			m = e
			t = i
		}
	}
	return t
}

func main() {
	var filename string
	var modelfile string
	flag.StringVar(&filename, "f", "", "input filename")
	flag.StringVar(&modelfile, "m", "mnist_model.tflite", "model filename")
	flag.Parse()

	var f *os.File
	var err error
	if filename == "" {
		f = os.Stdin
		filename = "stdin"
	} else {
		f, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	model := tflite.NewModelFromFile(modelfile)
	if model == nil {
		log.Println("cannot load model")
		return
	}
	defer model.Delete()

	interpreter := tflite.NewInterpreter(model, nil)
	defer interpreter.Delete()

	status := interpreter.AllocateTensors()
	if status != tflite.OK {
		log.Println("allocate failed")
		return
	}

	input := interpreter.GetInputTensor(0)
	resized := resize.Resize(28, 28, img, resize.NearestNeighbor)
	in := input.Float32s()
	for y := 0; y < 28; y++ {
		for x := 0; x < 28; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			in[y*28+x] = (float32(b) + float32(g) + float32(r)) / 3.0 / 65535.0
		}
	}
	status = interpreter.Invoke()
	if status != tflite.OK {
		log.Println("invoke failed")
		return
	}

	output := interpreter.GetOutputTensor(0)
	out := output.Float32s()
	fmt.Printf("%s is %d\n", filename, top(out))
}
