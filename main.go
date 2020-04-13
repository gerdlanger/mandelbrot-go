package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

type MyImage struct {
	Xsize, Ysize int
	XYrect       image.Rectangle
	buf          [][]color.Color
}

type MyPixel struct {
	X, Y, Iter int
}

type ImgRect struct {
	reMin, reMax float64
	imMin, imMax float64
}

func main() {

	numCPUs := runtime.NumCPU()
	numProcs := runtime.GOMAXPROCS(0)
	fmt.Printf("OS: %s Architecture: %s CPUs: %d Default Procs: %d\n", runtime.GOOS, runtime.GOARCH, numCPUs, numProcs)

	//const width, height = 3000, 2600
	width, height := 6000, 5200
	goXs, goYs := numCPUs, numCPUs

	examples := []ImgRect{
		ImgRect{-2.2, 0.8, -1.3, 1.3},
		ImgRect{0.435396403, 0.451687191, 0.367981352, 0.380210061},
		ImgRect{-0.37465401, -0.37332411, +0.659227668, +0.66020767},
		ImgRect{-1.816, -1.66, -0.078, 0.078},
		ImgRect{-0.098289376564, -0.098289374197, -0.889045279124, -0.88904527676},
		ImgRect{-0.098289375468734, -0.098289375467142, -0.889045277967886, -0.889045277966296},
		ImgRect{-0.867, -0.7, -0.234, 0},
		ImgRect{-0.8085, -0.78159, -0.19071, -0.1638},
		ImgRect{-0.80715, -0.7941, -0.18546, -0.17255},
		ImgRect{-0.8024335, -0.802384, -0.1776984, -0.1776492},
	}
	// Standard
	reMin, reMax := -2.2, 0.8
	imMin, imMax := -1.3, 1.3

	radius := 2.0
	radius2 := radius * radius
	maxIter := 1000

	args := os.Args[1:]
	argc := len(args)
	for idx := 0; idx < argc; idx++ {
		switch arg := args[idx]; arg {
		case "h", "?":
			fmt.Println("args:")
			fmt.Println("-c <num> => limit number of concurrent routines")
			fmt.Println("-p <num> => for NumProcs")
			fmt.Printf("-i <num> => max iteration (def=%d)\n", maxIter)
			fmt.Println("-s <size> => <size> is one of: (s) small, (xs) extra small, (tiny), (vga) =640x480, (hdmi)=1920x1080 ")
			fmt.Printf("-r <num> => select alternative image ranges (0-%d)\n", len(examples))
		case "-c":
			idx++
			num, err := extractArg(args[idx])
			if err == nil {
				switch num {
				case 1, 2, 3, 5:
					goXs, goYs = 1, num
				case 4, 6, 8, 10, 14:
					goXs, goYs = 2, num/2
				case 9, 15:
					goXs, goYs = 3, num/3
				case 12, 16, 32:
					goXs, goYs = 4, num/4
				case 64, 128:
					goXs, goYs = 8, num/8
				case 256:
					goXs, goYs = 16, 16
				case 100:
					goXs, goYs = 10, 10
				case 1000:
					goXs, goYs = 50, 20
				default:
					goXs, goYs = num, 1
				}
				fmt.Printf("c(%d) => x*y=%d*%d\n", num, goXs, goYs)
			}
		case "-p":
			idx++
			num, err := extractArg(args[idx])
			if err == nil {
				runtime.GOMAXPROCS(num)
				fmt.Printf("p(%d) => maxprocs=%d\n", num, num)
			}
		case "-i":
			idx++
			num, err := extractArg(args[idx])
			if err == nil {
				maxIter = num
				fmt.Printf("i(%d)\n", num)
			}
		case "-r":
			idx++
			num, err := strconv.Atoi(args[idx])
			if err == nil && num >= 0 && num < len(examples) {
				reMin, reMax = examples[num].reMin, examples[num].reMax
				imMin, imMax = examples[num].imMin, examples[num].imMax
				fmt.Printf("r(%d)\n", num)
			} else {
				fmt.Println("r out of range:", num)
			}
		case "-s":
			idx++
			switch size := args[idx]; size {
			case "s":
				width /= 2
				height /= 2
				fmt.Println("small pic")
			case "xs":
				width /= 4
				height /= 4
				fmt.Println("extra small pic")
			case "vga":
				width = 640
				height = 480
				fmt.Println("vga pic")
			case "hdmi":
				width = 1920
				height = 1080
				fmt.Println("hdmi pic")
			case "tiny":
				width = 320
				height = 200
				fmt.Println("tiny pic")
			default:
				fmt.Println("unknown size", size)
			}
		default:
			fmt.Println("unknown:", args, ". Use h for help")
		}
	}

	gos := goXs * goYs
	pixels := make(chan MyPixel)
	ready := make(chan MyPixel)

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	start := time.Now()

	initialGos := runtime.NumGoroutine()
	for y := 0; y < goYs; y++ {
		yPos1, ySize, minIm, maxIm := calculateRange(height, y, goYs, imMin, imMax)
		for x := 0; x < goXs; x++ {
			xPos1, xSize, minRe, maxRe := calculateRange(width, x, goXs, reMin, reMax)

			//fmt.Println(xPos1, yPos1, xPos2, yPos2, xSize, ySize, goHeight, goWidth)
			//fmt.Println(minRe, minIm, maxRe, maxIm, max2, xPos1, yPos1, xSize, ySize, maxIter)

			go Apfel(pixels, ready, minRe, minIm, maxRe, maxIm, radius2, float64(xPos1), float64(yPos1), float64(xSize), float64(ySize), maxIter)
			//go Apfel2(pixels, ready, minRe, minIm, maxRe, maxIm, max2, xPos1, yPos1, xSize, ySize, maxIter)
		}
	}
	fmt.Println("GoRoutines started:", runtime.NumGoroutine()-initialGos)

	running := gos
	for running > 0 {
		select {
		case pix := <-pixels:
			val := pix.Iter
			if val == maxIter {
				img.Set(pix.X, pix.Y, color.Black)
			} else {
				c := color.NRGBA{
					R: uint8((val) & 255),
					G: uint8((val) << 1 & 255),
					B: uint8((val) << 2 & 255),
					A: 255,
				}
				img.Set(pix.X, pix.Y, c)
			}
		case <-ready:
			running--
		}
	}

	duration := time.Now().Sub(start)
	fmt.Println("took: ", duration)

	txt := fmt.Sprintf("_%dx%d_%dx%d_%4.2f", width, height, goXs, goYs, duration.Seconds())

	save(img, "")
	save(img, txt)

	fmt.Println("done.")
}

func extractArg(arg string) (int, error) {
	num, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Println("invalid i:", arg)
	}
	return num, err
}

func save(img image.Image, txt string) error {
	name := "image" + txt + ".png"
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("saved:", name)
	}

	return err
}

func calculateRange(size, idx, total int, min, max float64) (pos, dist int, val1, val2 float64) {
	dSize := size / total
	delta := max - min

	pos = dSize * idx
	dist = dSize
	if pos+dist > size {
		dist = size - pos
	}
	val1 = min + delta*float64(pos)/float64(size)
	val2 = min + delta*float64(pos+dist)/float64(size)

	return
}

// Apfel caluclates an AppleMan
func Apfel(pixel, ready chan MyPixel, reMin, imMin, reMax, imMax, radius2, xoff, yoff, xs, ys float64, maxIter int) {

	imDelta := imMax - imMin
	reDelta := reMax - reMin
	//xsFloat, ysFloat := float64(xs), float64(ys)
	for y := 0.0; ys-y > 0.1; y++ {
		cIm := imMin + imDelta*y/ys
		for x := 0.0; xs-x > 0.1; x++ {
			cRe := reMin + reDelta*x/xs

			iter := Julia(cRe, cIm, cRe, cIm, radius2, maxIter)
			pixel <- MyPixel{int(x + xoff), int(y + yoff), iter}
		}
	}
	ready <- MyPixel{}
}

// Julia calculates single pixel by iterating (see https://de.wikipedia.org/wiki/Mandelbrot-Menge)
func Julia(x, y, xadd, yadd, radius2 float64, maxIter int) int {
	remain := maxIter
	xx := x * x
	yy := y * y
	xy := x * y
	len2 := xx + yy

	for ; len2 <= radius2 && remain > 0; remain-- {
		x = xx - yy + xadd
		y = xy + xy + yadd
		xx = x * x
		yy = y * y
		xy = x * y
		len2 = xx + yy
	}

	return maxIter - remain
}

// Apfel2 caluclates an AppleMan
func Apfel2(pixel, ready chan MyPixel, reMin, imMin, reMax, imMax, max2 float64, xoff, yoff, xs, ys, maxIter int) {

	imStep := (imMax - imMin) / float64(ys)
	reStep := (reMax - reMin) / float64(xs)
	cIm := imMin
	//xsFloat, ysFloat := float64(xs), float64(ys)
	for y := 0; y < ys; y++ {
		cRe := reMin
		for x := 0; x < xs; x++ {
			iter := Julia(cRe, cIm, cRe, cIm, max2, maxIter)
			pixel <- MyPixel{int(x + xoff), int(y + yoff), iter}
			cRe += reStep
		}
		cIm += imStep
	}
	ready <- MyPixel{}
}
