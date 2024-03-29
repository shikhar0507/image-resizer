package filters

import (
	"image"
	"image/color"
	"math"
)

func init() {
}

type BrightnessOptions struct {
	Factor float64
}
type BlurOptions struct {
	Radius int
}

func BrightnessAdjust(dec image.Image, width float64, height float64, bo BrightnessOptions) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	for x := 0; x < int(width); x++ {
		for y := 0; y < int(height); y++ {
			r, g, b, a := dec.At(x, y).RGBA()
			newImage.Set(x, y, color.RGBA{R: getAdjustedPixel(r, bo.Factor, 1), G: getAdjustedPixel(g, bo.Factor, 1), B: getAdjustedPixel(b, bo.Factor, 1), A: getAdjustedPixel(a, bo.Factor, 1)})
		}
	}

	return newImage
}

func getAdjustedPixel(value uint32, factor float64, alpha float64) uint8 {
	px := uint8(value / 257)
	adjusted := (alpha * float64(px)) + factor
	if adjusted < 0 {
		return 0
	}
	if adjusted > 255 {
		return 255
	}
	return uint8(adjusted)
}

func BlackAndWhite(dec image.Image, width float64, height float64) image.Image {
	newImage := image.NewGray(image.Rect(0, 0, int(width), int(height)))
	grayScaled := Grayscale(dec, width, height)

	for x := 0; x < int(width); x++ {
		for y := 0; y < int(height); y++ {
			grayValue, _, _, _ := grayScaled.At(x, y).RGBA()
			uintGray := uint8(grayValue / 257)
			if 255-uintGray < uintGray {
				uintGray = 255
			} else {
				uintGray = 0

			}
			newImage.SetGray(x, y, color.Gray{Y: uintGray})
		}
	}

	return newImage
}

func BlackAndWhiteSmooth(dec image.Image, width float64, height float64) image.Image {
	newImage := image.NewGray(image.Rect(0, 0, int(width), int(height)))
	transformed := BlackAndWhite(dec, width, height)

	for x := 0; x < int(width); x++ {
		for y := 0; y < int(height); y++ {
			gr := uint8(255)
			blackCount, whiteCount := 0, 0
			for i := -3; i <= 3; i++ {
				for j := -3; j <= 3; j++ {
					val, _, _, _ := transformed.At(x, y).RGBA()
					if val == 0 {
						blackCount++
					} else {
						whiteCount++
					}
				}
			}
			if blackCount > whiteCount {
				gr = 0
			}
			newImage.SetGray(x, y, color.Gray{Y: gr})
		}
	}
	return newImage
}

func Sepia(dec image.Image, width, height float64) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	for x := 0; x < int(width); x++ {
		for y := 0; y < int(height); y++ {
			r, g, b, _ := dec.At(x, y).RGBA()

			r8, g8, b8 := float64(uint8(r/257)), float64(uint8(g/257)), float64(uint8(b/257))

			nr := r8*0.393 + g8*0.769 + b8*0.189
			if nr > 255 {
				nr = 255
			}
			ng := r8*0.349 + g8*0.686 + b8*0.168
			if ng > 255 {
				ng = 255
			}
			nb := r8*0.272 + g8*0.534 + b8*0.131
			if nb > 255 {
				nb = 255
			}
			newImage.Set(x, y, color.RGBA{
				R: uint8(nr),
				G: uint8(ng),
				B: uint8(nb),
			})
		}

	}

	return newImage
}

func Blur(dec image.Image, width float64, height float64, bo BlurOptions) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	blurOff := (bo.Radius - 1) / 2

	for x := 0; x < int(width); x++ {
		for y := 0; y < int(height); y++ {
			totalR, totalG, totalB, totalA := 0, 0, 0, 0
			size := 0
			for i := -blurOff; i <= blurOff; i++ {
				for j := -blurOff; j <= blurOff; j++ {
					size++
					pxR, pxG, pxB, pxA := dec.At(x+i, y+j).RGBA()
					totalR += int(pxR)
					totalG += int(pxG)
					totalB += int(pxB)
					totalA += int(pxA)
				}
			}

			newImage.SetRGBA(x, y, color.RGBA{R: uint8((totalR / size) / 257), G: uint8((totalG / size) / 257), B: uint8((totalB / size) / 257), A: uint8((totalA / size) / 257)})
		}
	}
	return newImage
}

func Grayscale(dec image.Image, width float64, height float64) image.Image {
	newImage := image.NewGray(image.Rect(0, 0, int(width), int(height)))
	sourceWidth := float64(dec.Bounds().Size().X)
	sourceHeight := float64(dec.Bounds().Size().Y)
	for i := float64(0); i < width; i++ {
		x := int(math.Round((i / width) * sourceWidth))
		for j := float64(0); j < height; j++ {
			y := int(math.Round((j / height) * sourceHeight))
			col := dec.At(x, y)
			r, g, b, _ := col.RGBA()

			R, G, B := float64(r)*0.299, float64(g)*0.587, float64(b)*0.114
			tot := uint8((R + G + B) / 257)
			//fmt.Println(R,G,B)
			//uintR ,uintG,uintb := uint8(r/257),uint8(g/257),uint8(b/257)
			newImage.SetGray(int(i), int(j), color.Gray{Y: tot})
		}
	}
	return newImage
}

func Negative(dec image.Image, width float64, height float64) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for x := 0; x < int(width); x++ {
		for y := 0; y < int(height); y++ {
			r, g, b, _ := dec.At(x, y).RGBA()
			newR, newG, newB := 255-getAdjustedPixel(r, 0, 1), 255-getAdjustedPixel(g, 0, 1), 255-getAdjustedPixel(b, 0, 1)
			newImage.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB})
		}
	}
	return newImage
}

func Postive(dec image.Image, width float64, height float64) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for x := 0; x < int(width); x++ {
		for y := 0; y < int(height); y++ {
			r, g, b, _ := dec.At(x, y).RGBA()
			newR, newG, newB := 255-getAdjustedPixel(r, 0, 1), 255-getAdjustedPixel(g, 0, 1), 255-getAdjustedPixel(b, 0, 1)
			newImage.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB})
		}
	}
	return newImage
}
