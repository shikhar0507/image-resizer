package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"math"
	"os"
	//"os"

	//"io"
	//"io/ioutil"
	"log"
	"net/http"
	//"os"
	"strconv"
)

//type Extension struct {
//	jpg int
//	png int
//	jpeg int
//}
type Properties struct {
	url string
	width ,height float64
	//ext Extension
}
//var extension Extension
var property Properties


func main() {

	//extension = Extension{jpeg: 1,png: 2,jpg: 3}
	http.HandleFunc("/",handleRequest)
	http.HandleFunc("/favicon.ico",handleFavicon)

	log.Fatal(http.ListenAndServe(":8080",nil))
}

func handleFavicon(w http.ResponseWriter , r *http.Request) {

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
	query := r.URL.Query()
	url,width,height := query.Get("url"),query.Get("width"),query.Get("height")

	if url == "" {
		http.Error(w,"url is not provided",http.StatusBadRequest)
		return
	}

	if width == "" {
		http.Error(w,"width is not provided",http.StatusBadRequest)
	}
	if height == "" {
		http.Error(w,"height is not provided",http.StatusBadRequest)
	}
	wd,werr := strconv.Atoi(width)
	hg, herr := strconv.Atoi(height)

	if werr !=nil {
		http.Error(w,"error parsing the width",http.StatusInternalServerError)
		return
	}
	if herr != nil {
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

	property = Properties{url: url,width: float64(wd),height: float64(hg)}
	loadImageFromUrl(property)
}

func loadImageFromUrl(prop Properties) (image.Image, error) {
	resp , err := http.Get(prop.url)
	if err != nil {
		fmt.Println("error loading image.jpeg",err)
		return nil,err
	}
	defer resp.Body.Close()


	dec ,err := jpeg.Decode(resp.Body)
	newImage := image.NewRGBA(image.Rect(0,0,int(prop.width),int(prop.height)))


	if err != nil {
		log.Fatal("error decoding jpeg",err)
	}
	var i float64 = 0.0

	sourceWidth := float64(dec.Bounds().Size().X)
	sourceHeight := float64(dec.Bounds().Size().Y)
	for  ; i< prop.width;i++ {
		x :=  int(math.Round((i/prop.width)* sourceWidth))
		for j := float64(0);j<prop.height;j++ {
			y := int(math.Round((j/prop.height) * sourceHeight))
			col := dec.At(x,y)
			r,g,b,a := col.RGBA()

			newImage.SetRGBA(int(i),int(j),color.RGBA{R: uint8(r/257),G: uint8(g/257),B: uint8(b/257),A: uint8(a/257)})
		}
	}
	fmt.Println("new image bounds", newImage.Bounds().Size())
	//draw.Draw(newImage,newImage.Bounds(),newImage,image.Pt(0,0),draw.Src)

	f, err := os.Create("i.jpg")
	if err != nil {
		log.Fatal(err)
	}

	//var opt Options
	//opt = Options{Quality: 75}
	jpeg.Encode(f,newImage,&jpeg.Options{Quality: 100})
	fmt.Println("done")
	return nil,nil
}
