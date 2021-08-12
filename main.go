package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image-resizer/filters"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
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

	fmt.Println("starting image resizing")

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/favicon.ico", handleFavicon)
	log.Fatal(http.ListenAndServe(":5001", nil))
	//img := <-ch
	//close(ch)

	/*
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error", http.StatusBadRequest)
			return
		}
	*/

}

func handleFavicon(w http.ResponseWriter, r *http.Request) {

}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	fmt.Println("request incoming")
	query := r.URL.Query()
	urlStr, width, height, filter, value := query.Get("url"), query.Get("width"), query.Get("height"), query.Get("filter"), query.Get("value")

	if urlStr == "" {
		http.Error(w, "url is not provided", http.StatusBadRequest)
		return
	}

	wd, werr := strconv.Atoi(width)
	hg, herr := strconv.Atoi(height)
	var intVal int

	if width == "" && height != "" {
		http.Error(w, "error parsing the width", http.StatusInternalServerError)
		return
	}
	if height == "" && width != "" {
		http.Error(w, "error parsing the height", http.StatusInternalServerError)
		return
	}

	if width != "" && werr != nil {
		http.Error(w, "error parsing the width", http.StatusInternalServerError)
		return
	}

	if height != "" && herr != nil {
		http.Error(w, "error parsing the height", http.StatusInternalServerError)
		return
	}

	if wd <= 0 {
		http.Error(w, "width cannot be less than or equal to 0", http.StatusBadRequest)
		return
	}
	if hg <= 0 {
		http.Error(w, "height cannot be less than or equal to 0", http.StatusBadRequest)
		return

	}
	if value != "" {
		val, convErr := strconv.Atoi(value)
		if convErr != nil {
			intVal = 0
		} else {
			intVal = val
		}

	}
	//parsedUrl,err :=  url.Parse(urlStr)
	fmt.Println(intVal)
	property = Properties{url: r.URL, width: float64(wd), height: float64(hg), filter: filter, imageUrl: urlStr, value: intVal}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "max-age=3600")
	//	ch := make(chan image.Image)

	var ch chan ChanStr = make(chan ChanStr)

	go loadImageFromUrl(property, w, ch)
	for result := range ch {
		if result.img == nil {
			log.Fatal(result.img)
			return
		}
		buffer := new(bytes.Buffer)
		fmt.Println("encoding image")

		encodeErr := jpeg.Encode(buffer, result.img, &jpeg.Options{Quality: 100})

		//encodeErr := jpeg.Encode(buffer, img, &jpeg.Options{Quality: 100})
		if encodeErr != nil {
			http.Error(result.w, "error", http.StatusInternalServerError)
			return
		}

		_, writeErr := result.w.Write(buffer.Bytes())
		if writeErr != nil {
			log.Println(writeErr)
			return
		}
	}

	//close(ch)
}

func getImage(imageUrl string) (io.Reader, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println("error loading image", err)
		return nil, err
	}
	respBody := resp.Body
	dat, err := ioutil.ReadAll(respBody)
	return bytes.NewReader(dat), nil

}

func loadImageFromUrl(prop Properties, w http.ResponseWriter, ch chan ChanStr) {

	respBody, err := getImage(prop.imageUrl)
	if err != nil {
		ch <- ChanStr{}
		return
	}
	decodedImage, _, err := image.Decode(respBody)

	//fmt.Println("type",imageType)

	if err != nil {
		ch <- ChanStr{}
		return
	}
	fmt.Println("decoded image")
	if decodedImage == nil {
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

	ch <- ChanStr{img: modifiedImage, w: w}
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
