package utils

import (
	"fmt"
	"image"
)

func init() {
	fmt.Println("init utils")

}

func ConvertToHSL(dec image.Image, width, height int) image.Image {
	hslImage := image.NewRGBA(image.Rect(0, 0, width, height))

	return hslImage
}

func ConvertToRGB(dec image.Image, width, height int) image.Image {
	hslImage := image.NewRGBA(image.Rect(0, 0, width, height))

	return hslImage
}
