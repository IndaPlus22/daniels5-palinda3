// Stefan Nilsson 2013-02-27

// This program creates pictures of Julia sets (en.wikipedia.org/wiki/Julia_set).

//Original run time: 11.29s user 0.13s system 99% cpu 11.519 total
//Improvement #1 seperating each picture generation into seperate routines: 11.58s user 0.17s system 144% cpu 8.106 total
//Imporovment #2 Seperating the the pixel generation process to different routines
//This is done by letting each routine handle each part of the image seperatly
//Lets say i use 2 routines then the first routine would color pixel 0, 2, 4, 6, 8 and so on
//The other would handle pixel 1, 3, 5, 7, ... and so on.
//If i use 10 routines the speed gets down to 2.236 seconds. Me after seeing that result: https://www.meme-arsenal.com/memes/6a6268068ad5ef424ab9a50cb0d02f8a.jpg

package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
	"strconv"
	"sync"
)

type ComplexFunc func(complex128) complex128

var Funcs []ComplexFunc = []ComplexFunc{
	func(z complex128) complex128 { return z*z - 0.61803398875 },
	func(z complex128) complex128 { return z*z + complex(0, 1) },
	func(z complex128) complex128 { return z*z + complex(-0.835, -0.2321) },
	func(z complex128) complex128 { return z*z + complex(0.45, 0.1428) },
	func(z complex128) complex128 { return z*z*z + 0.400 },
	func(z complex128) complex128 { return cmplx.Exp(z*z*z) - 0.621 },
	func(z complex128) complex128 { return (z*z+z)/cmplx.Log(z) + complex(0.268, 0.060) },
	func(z complex128) complex128 { return cmplx.Sqrt(cmplx.Sinh(z*z)) + complex(0.065, 0.122) },
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(len(Funcs))
	for n, fn := range Funcs {

		go CreatePng("picture-"+strconv.Itoa(n)+".png", fn, 1024, wg)
	}
	wg.Wait()

}

// CreatePng creates a PNG picture file with a Julia image of size n x n.
func CreatePng(filename string, f ComplexFunc, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	err = png.Encode(file, Julia(f, n))
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

// Julia returns an image of size n x n of the Julia set for f.
func Julia(f ComplexFunc, n int) image.Image {
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	wg := new(sync.WaitGroup)
	numOfGenerators := 10
	for i := 0; i < numOfGenerators; i++ {
		wg.Add(1)
		go severalPartPixelGen(f, numOfGenerators, i, n, img, wg)
	}
	wg.Wait()
	return img
}

func severalPartPixelGen(f ComplexFunc, step, startPos, n int, img *image.RGBA, wg *sync.WaitGroup) {
	defer wg.Done()
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	s := float64(n / 4)

	for i := bounds.Min.X + startPos; i < bounds.Max.X; i += step {
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			n := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
			r := uint8(0)
			g := uint8(0)
			b := uint8(n % 32 * 8)
			img.Set(i, j, color.RGBA{r, g, b, 255})
		}
	}
}

// Iterate sets z_0 = z, and repeatedly computes z_n = f(z_{n-1}), n â‰¥ 1,
// until |z_n| > 2  or n = max and returns this n.
func Iterate(f ComplexFunc, z complex128, max int) (n int) {
	for ; n < max; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
		z = f(z)
	}
	return
}
