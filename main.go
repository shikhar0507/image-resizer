package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"log"
	"net/http"
	"strconv"
	"image-resizer/filters"
	//"github.com/shikhar0507/image-resizer"
)


type Properties struct {
	url string
	width ,height float64
	filter string
	//ext Extension
}

var property Properties


func main() {

	http.HandleFunc("/",handleRequest)
	http.HandleFunc("/favicon.ico",handleFavicon)
	log.Fatal(http.ListenAndServe(":8080",nil))


}

func handleFavicon(w http.ResponseWriter , r *http.Request) {

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
	query := r.URL.Query()
	url,width,height,filter := query.Get("url"),query.Get("width"),query.Get("height"),query.Get("filter")

	if url == "" {
		http.Error(w,"url is not provided",http.StatusBadRequest)
		return
	}


	wd,werr := strconv.Atoi(width)
	hg, herr := strconv.Atoi(height)

	if width == "" && height != "" {
		http.Error(w,"error parsing the width",http.StatusInternalServerError)
		return
	}
	if height == "" && width != "" {
		http.Error(w,"error parsing the height",http.StatusInternalServerError)
		return
	}


	if width != "" && werr !=nil {
		http.Error(w,"error parsing the width",http.StatusInternalServerError)
		return
	}

	if height != "" && herr != nil {
		http.Error(w,"error parsing the height",http.StatusInternalServerError)
		return
	}
	
	if wd <= 0 {
		http.Error(w,"width cannot be less than or equal to 0",http.StatusBadRequest)
		return
	}
	if hg <= 0 {
		http.Error(w,"height cannot be less than or equal to 0",http.StatusBadRequest)
		return
	}

	property = Properties{url: url,width: float64(wd),height: float64(hg),filter:filter}
	loadImageFromUrl(property)
}



func loadImageFromUrl(prop Properties) (image.Image, error) {
	resp , err := http.Get(prop.url)
	if err != nil {
		fmt.Println("error loading image.jpeg",err)
		return nil,err
	}
	defer resp.Body.Close()
	respBody := resp.Body
	decodedImage,_ ,err := image.Decode(respBody)
	//fmt.Println("type",imageType)

	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	if decodedImage == nil {
		fmt.Println("null image")
		return nil,nil
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

	var baseImage image.Image
	if imageBound.X != int(prop.width) || imageBound.Y != int(prop.height) {
		fmt.Println("start resizing")
	   baseImage = resize(decodedImage,prop.width,prop.height, float64(imageBound.X),float64(imageBound.Y))
	}

	var modifiedImage image.Image

	switch prop.filter {
	       case "sepia":
	       	    modifiedImage = filters.Sepia(baseImage,sourceWidth,sourceHeight)
	       case "blackAndWhite":
	       	    modifiedImage = filters.BlackAndWhite(baseImage,sourceWidth,sourceHeight)
	       case "blackAndWhiteSmooth":
	       	    modifiedImage = filters.BlackAndWhiteSmooth(baseImage,sourceWidth,sourceHeight)
	       case "grayscale":
	       	    modifiedImage = filters.Grayscale(baseImage,sourceWidth,sourceHeight)
	       case "brightness":
	       	    bo := filters.BrightnessOptions{Factor:150}
	       	    modifiedImage = filters.BrightnessAdjust(baseImage,sourceWidth,sourceHeight,bo)
	      case "blur":
	      	   bo := filters.BlurOptions{Radius:-3}
		   modifiedImage = filters.Blur(baseImage,sourceWidth,sourceHeight,bo)
	     default:
		   modifiedImage = baseImage
	}

	

	f, err := os.Create("output.jpg")
	if err != nil {
		log.Fatal(err)
	}


	jpeg.Encode(f,modifiedImage,&jpeg.Options{Quality: 100})
	fmt.Println("done")
	return nil,nil
}

func resize(dec image.Image,width float64,height float64,sourceWidth float64,sourceHeight float64) image.Image {
     newImage := image.NewRGBA(image.Rect(0,0,int(width),int(height)))
     for i:= float64(0) ; i< width;i++ {
     	 x :=  int(math.Round((i/width)* sourceWidth))
  	   for j := float64(0);j<height;j++ {
               y := int(math.Round((j/height) * sourceHeight))
               col := dec.At(x,y)
               r,g,b,a := col.RGBA()
	       uintR ,uintG,uintB,uintA := uint8(r/257),uint8(g/257),uint8(b/257),uint8(a/257)
	       newImage.SetRGBA(int(i),int(j),color.RGBA{R: uintR,G: uintG,B: uintB,A:uintA})
  	       }
	  }
     return newImage

}