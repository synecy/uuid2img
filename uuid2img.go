package uuid2img

import (
	"fmt"
	"math/big"
	"strings"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Pixel struct {
	R int
	G int
	B int
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r >> 8), int(g >> 8), int(b >> 8)}
}

func saveImage(fileName string, pixels [][]Pixel) bool {
	width := len(pixels[0])
	height := len(pixels)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8(pixels[y][x].R),
				G: uint8(pixels[y][x].G),
				B: uint8(pixels[y][x].B),
				A: 255,
			})
		}
	}

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		fmt.Println(err)
		return false
	}

	if err := f.Close(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func newImage(width int, height int) [][]Pixel {
	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var pixelRow []Pixel
		for x := 0; x < width; x++ {
			pixel := rgbaToPixel(255, 255, 255, 255)
			pixelRow = append(pixelRow, pixel)
		}
		pixels = append(pixels, pixelRow)
	}
	return pixels
}

func copyMatrixToImage(pixels [][]Pixel, matrix [][]int) [][]Pixel {
	for y := 0; y < len(matrix); y++ {
		for x := 0; x < len(matrix[0]); x++ {
			pixels[y][x].R = matrix[y][x]
			if x+1 < len(matrix[0]) {
				pixels[y][x].G = matrix[y][x+1]
			} else {
				pixels[y][x].G = matrix[y][x-1]
			}
			if x+2 < len(matrix[0]) {
				pixels[y][x].B = matrix[y][x+2]
			} else {
				pixels[y][x].B = matrix[y][x-2]
			}
		}
	}
	return pixels
}

func resizeImage(pixels [][]Pixel, multiplier int) [][]Pixel {
	width := len(pixels) * multiplier
	height := len(pixels[0]) * multiplier
	newImg := newImage(width, height)
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[0]); x++ {
			for yi := 0; yi < multiplier; yi++ {
				for xi := 0; xi < multiplier; xi++ {
					newImg[(y*multiplier)+yi][(x*multiplier)+xi].R = pixels[y][x].R
					newImg[(y*multiplier)+yi][(x*multiplier)+xi].G = pixels[y][x].G
					newImg[(y*multiplier)+yi][(x*multiplier)+xi].B = pixels[y][x].B
				}
			}
		}
	}
	return newImg
}

func GenerateFile(uid string, filepath string) bool {
	uidPairs := make([]string, 0)
	trimmedUid := strings.ReplaceAll(uid, "-", "")
	uidTokens := strings.SplitN(trimmedUid, "", len(trimmedUid))
	for i := 0; i < len(uidTokens)-1; i += 2 {
		uidPairs = append(uidPairs, uidTokens[i]+uidTokens[i+1])
	}
	uidMatrix := make([][]int, 0)
	rowArr := make([]int, 0)
	for index, hexNum := range uidPairs {
		n := new(big.Int)
		colorCode, _ := n.SetString(hexNum, 16)
		rowArr = append(rowArr, int(colorCode.Int64()))
		if index%4 == 3 {
			uidMatrix = append(uidMatrix, rowArr)
			rowArr = nil
		}
	}
	newImg := newImage(4, 4)
	newImg = copyMatrixToImage(newImg, uidMatrix)
	newImg = resizeImage(newImg, 32)
	return saveImage(filepath, newImg)
}
