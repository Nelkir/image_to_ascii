package main

import (
	//	"encoding/base64"

	"flag"
	"fmt"
	"image"
	"log"
	"time"

	//	"strings"
	"os"
	//	"reflect"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Uncomment these
	// two lines to also understand GIF and PNG images:
	// _ "image/gif"
	// _ "image/png"

	"image/jpeg"
	_ "image/jpeg"

	"github.com/nfnt/resize"
)

var (
	path_to_image        = flag.String("i", "", "set path image")
	width                = flag.Int("w", 0, "Set with of resulting image")
	height               = flag.Int("h", 0, "Set height of resulting image")
	debug                = flag.Bool("d", false, "Set debug for timings. Write to console. Image still writes direct to stdin")
	as_char              = flag.Bool("c", false, "Print string as array of chars")
	negative             = flag.Bool("n", false, "Reverse brighness detection")
	create_resized_image = flag.String("r", "", "Create resized copy of image")
)

func main() {
	start := time.Now()
	scale := "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
	flag.Parse()
	if *path_to_image == "" {
		fmt.Println("set path with -i: ", flag.ErrHelp)
		return
	}
	since := time.Since(start)
	args_time := fmt.Sprintf("Args took = %v\n", since)

	start = time.Now()
	reader, err := os.Open(*path_to_image)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	since = time.Since(start)
	read_time := fmt.Sprintf("Read image took = %v\n", since)

	start = time.Now()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	since = time.Since(start)
	decode_time := fmt.Sprintf("Decode image took = %v\n", since)

	start = time.Now()
	resized_image := resize.Resize(uint(*width), uint(*height), m, resize.Bicubic)
	resized_bounds := resized_image.Bounds()
	since = time.Since(start)
	resize_time := fmt.Sprintf("Resize image took = %v\n", since)

	start = time.Now()
	screen := make([][]byte, resized_bounds.Max.Y)
	for i := 0; i < resized_bounds.Max.Y; i++ {
		screen[i] = make([]byte, resized_bounds.Max.X)
	}
	for i := 0; i < resized_bounds.Max.X; i++ {
		for j := 0; j < resized_bounds.Max.Y; j++ {
			r, g, b, _ := resized_image.At(i, j).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			l := (0.2126*float32(r) + 0.7152*float32(g) + 0.0722*float32(b))
			var lum float32
			if *negative {
				lum = 69 + (-l*float32(69))/float32(255)
			} else {
				lum = (l * float32(69)) / float32(255)
			}
			screen[j][i] = scale[uint8(lum)]
		}
	}
	since = time.Since(start)
	prepare_time := fmt.Sprintf("Prepare screen buff took = %v\n", since)

	if *create_resized_image != "" {
		resized_image_file, err := os.OpenFile(*create_resized_image, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Printf("Can't create resized image: %s\n", err)
		}
		jpeg.Encode(resized_image_file, resized_image, nil)
	}

	start = time.Now()
	for i := 0; i < resized_bounds.Max.Y; i++ {
		if *as_char {
			fmt.Fprintf(os.Stdin, "%c\n", screen[i])
		} else {
			fmt.Fprintf(os.Stdin, "%s\n", screen[i])
		}
	}
	since = time.Since(start)
	if *debug {
		fmt.Print(args_time, read_time, decode_time, resize_time, prepare_time, fmt.Sprintf("Print screen buff took = %v\n", since))
	}
}
