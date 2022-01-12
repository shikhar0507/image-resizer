package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image-resizer/filters"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Properties struct {
	url           *url.URL
	imageUrl      string
	width, height float64
	filter        string
	value         int
	//ext Extension
}

type ChanStr struct {
	img image.Image
	w   http.ResponseWriter
}

var property Properties
var ctx = context.Background()

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, evt events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	head := map[string]string{
		"Content-Type":  "image/jpeg",
		"Cache-Control": "max-age=63072000",
		"Accept-Ranges": "bytes",
	}

	req, err := handleRequest(evt)
	if err != nil {

		return events.APIGatewayProxyResponse{Headers: head, StatusCode: 500, Body: err.Error()}, err
	}

	return events.APIGatewayProxyResponse{Body: req, StatusCode: 200, Headers: head, IsBase64Encoded: true}, nil
}

func handleRequest(evt events.APIGatewayProxyRequest) (string, error) {
	var urlStr, width, height, filter, value string
	for k, v := range evt.QueryStringParameters {
		switch k {
		case "url":
			urlStr = v
		case "width":
			width = v
		case "height":
			height = v
		case "filter":
			filter = v
		case "value":
			value = v
		}

	}

	wd, werr := strconv.Atoi(width)
	hg, herr := strconv.Atoi(height)
	var intVal int

	if werr != nil {
		return "", errors.New("error parsing the width")

	}

	if herr != nil {
		return "", errors.New("error parsing the height")
	}

	if wd <= 0 {
		return "", errors.New("width cannot be less than or equal to 0")

	}
	if hg <= 0 {
		return "", errors.New("height cannot be less than or equal to 0")

	}
	if value != "" {
		val, convErr := strconv.Atoi(value)
		if convErr != nil {
			intVal = 0
		} else {
			intVal = val
		}

	}

	fmt.Println(intVal)
	ogUrl, _ := url.Parse(urlStr)
	property = Properties{url: ogUrl, width: float64(wd), height: float64(hg), filter: filter, imageUrl: urlStr, value: intVal}

	var ch chan ChanStr = make(chan ChanStr)

	go loadImageFromUrl(property, ch)
	result := <-ch

	if result.img == nil {
		return "", errors.New("image not found")
	}
	buffer := new(bytes.Buffer)

	encodeErr := jpeg.Encode(buffer, result.img, &jpeg.Options{Quality: 100})
	if encodeErr != nil {
		return "", encodeErr

	}

	return base64.StdEncoding.EncodeToString([]byte(buffer.Bytes())), nil
}

func getImage(imageUrl string) (io.Reader, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		return nil, err
	}
	//defer resp.Body.Close()
	return resp.Body, nil

}

func loadImageFromUrl(prop Properties, ch chan ChanStr) {

	respBody, err := getImage(prop.imageUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		ch <- ChanStr{}
		return
	}
	decodedImage, err := jpeg.Decode(respBody)

	//fmt.Println("type",imageType)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		ch <- ChanStr{}
		return
	}
	fmt.Println("decoded image")
	if decodedImage == nil {
		fmt.Fprintf(os.Stderr, "%s", "decoded image is nil")
		ch <- ChanStr{}
		return
	}

	imageBound := decodedImage.Bounds().Size()
	sourceWidth := float64(imageBound.X)
	if prop.width != 0 {
		sourceWidth = prop.width
	}

	sourceHeight := float64(imageBound.Y)
	if prop.height != 0 {
		sourceHeight = prop.height
	}

	baseImage := decodedImage
	if imageBound.X != int(prop.width) || imageBound.Y != int(prop.height) {
		fmt.Println("start resizing")
		baseImage = resize(decodedImage, prop.width, prop.height, float64(imageBound.X), float64(imageBound.Y))
	}

	var modifiedImage image.Image
	switch prop.filter {
	case "sepia":
		modifiedImage = filters.Sepia(baseImage, sourceWidth, sourceHeight)
	case "blackAndWhite":
		modifiedImage = filters.BlackAndWhite(baseImage, sourceWidth, sourceHeight)
	case "blackAndWhiteSmooth":
		modifiedImage = filters.BlackAndWhiteSmooth(baseImage, sourceWidth, sourceHeight)
	case "grayscale":
		modifiedImage = filters.Grayscale(baseImage, sourceWidth, sourceHeight)
	case "brightness":

		bo := filters.BrightnessOptions{Factor: float64(prop.value)}
		modifiedImage = filters.BrightnessAdjust(baseImage, sourceWidth, sourceHeight, bo)
	case "blur":
		bo := filters.BlurOptions{Radius: prop.value}
		modifiedImage = filters.Blur(baseImage, sourceWidth, sourceHeight, bo)
	case "color":
		modifiedImage = filters.Color(baseImage, sourceWidth, sourceHeight)

	default:
		modifiedImage = baseImage
	}

	ch <- ChanStr{img: modifiedImage}
	close(ch)
}

func resize(dec image.Image, width float64, height float64, sourceWidth float64, sourceHeight float64) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for i := float64(0); i < width; i++ {
		x := int(math.Round((i / width) * sourceWidth))
		for j := float64(0); j < height; j++ {
			y := int(math.Round((j / height) * sourceHeight))
			col := dec.At(x, y)
			r, g, b, a := col.RGBA()
			uintR, uintG, uintB, uintA := uint8(r/257), uint8(g/257), uint8(b/257), uint8(a/257)
			newImage.SetRGBA(int(i), int(j), color.RGBA{R: uintR, G: uintG, B: uintB, A: uintA})
		}
	}

	return newImage

}
