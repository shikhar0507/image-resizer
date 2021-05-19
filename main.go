package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
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
	"time"
)

type Properties struct {
	url           *url.URL
	imageUrl string
	width, height float64
	filter        string
	value int
	//ext Extension
}

var property Properties
var ctx = context.Background()
var rdb *redis.Client

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	fmt.Println("starting image resizing")	
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/favicon.ico", handleFavicon)
	log.Fatal(http.ListenAndServe(":5001", nil))

}

func handleFavicon(w http.ResponseWriter, r *http.Request) {

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
       fmt.Println("request incoming")
        query := r.URL.Query()
	urlStr, width, height, filter,value := query.Get("url"), query.Get("width"), query.Get("height"), query.Get("filter"), query.Get("value")
	

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
	   val,convErr := strconv.Atoi(value)
	   if convErr != nil {
	      intVal = 0
	   }else {intVal = val}
	   
	}
	//parsedUrl,err :=  url.Parse(urlStr)
	fmt.Println(intVal)
	property = Properties{url: r.URL, width: float64(wd), height: float64(hg), filter: filter,imageUrl: urlStr,value: intVal}
	
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "max-age=100")

	imageHash := sha256.Sum256([]byte(r.URL.RawQuery))

	savedImage, err := rdb.Get(ctx, fmt.Sprintf("%x", imageHash)).Result()
	if err == redis.Nil {
		img, err := loadImageFromUrl(property)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error", http.StatusBadRequest)
			return
		}
		if img == nil {
			log.Fatal(img)
			return
		}
		buffer := new(bytes.Buffer)
		fmt.Println("encoding image")
		encodeErr := jpeg.Encode(buffer, img, &jpeg.Options{Quality: 100})
		if encodeErr != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		rdb.Set(ctx, fmt.Sprintf("%x", imageHash), buffer.Bytes(), 20*time.Second)
		_, writeErr := w.Write(buffer.Bytes())
		if writeErr != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("image from redis found")
		w.Write([]byte(savedImage))
	}
}

func getImage(imageUrl string) (io.Reader, error) {
	key := fmt.Sprintf("%x", sha256.Sum256([]byte(imageUrl)))
	savedImageResponse, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("image not found in redis with key", key)
		fmt.Println("downloading image", imageUrl)
		resp, err := http.Get(imageUrl)
		if err != nil {
			fmt.Println("error loading image", err)
			return nil, err
		}
		//defer resp.Body.Close()

		respBody := resp.Body
		fmt.Println("storing image in redis with key", key)
		//r := bufio.NewReader(respBody)
		dat, err := ioutil.ReadAll(respBody)

		//m, err := json.Marshal(respBody)
		//fmt.Println("marshalled length",)
		s, e := rdb.Set(ctx, key, dat, 20*time.Second).Result()
		if e != nil {
			fmt.Println(e)
		}
		fmt.Println(s)
		fmt.Println("sending image")

		//buff := new(bytes.Buffer)
		return bytes.NewReader(dat) , nil
	}
	if err != nil {
		fmt.Println("reading err", err)
		return nil, err
	}
	fmt.Println("sending saved imaged")
	//var ma map[string]string
	////json.Unmarshal([]byte(savedImageResponse),ma)
	//fmt.Println(ma)
	//fmt.Println(savedImageResponse)
	return bytes.NewReader([]byte(savedImageResponse)) , nil
	//return bytes.NewReader([]byte(ma)), nil
}

func loadImageFromUrl(prop Properties) (image.Image, error) {

	respBody, err := getImage(prop.imageUrl)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	decodedImage, _, err := image.Decode(respBody)

	//fmt.Println("type",imageType)

	if err != nil {
		fmt.Println("decode err", err)
		return nil, err
	}
	fmt.Println("decoded image")
	if decodedImage == nil {
		fmt.Println("null image")
		return nil, errors.New("null image")
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
	default:
		modifiedImage = baseImage
	}

	

	return modifiedImage, nil

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
