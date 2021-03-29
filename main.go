package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
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
	filter string
	//ext Extension
}
=
//var extension Extension
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
	url,width,height := query.Get("url"),query.Get("width"),query.Get("height")

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

	//property = Properties{url: url,width: float64(wd),height: float64(hg)}
	loadImageFromUrl(property)
}



func loadImageFromUrl(prop Properties) (image.Image, error) {
	resp , err := http.Get(prop.url)
	if err != nil {
		fmt.Println("error loading image.jpeg",err)
		return nil,err
	}
	defer resp.Body.Close()
	dec,imageType, err := image.Decode(resp.Body)
	var decodedImage image.Image
	var decodedError error
	switch imageType {
	case "png":
		decodedImage, decodedError = png.Decode(resp.Body)
	case "jpg":
		decodedImage, decodedError = jpeg.Decode(resp.Body)
	}
	if decodedError != nil {
		return nil,decodedError
	}


	sourceWidth := float64(dec.Bounds().Size().X)
	if prop.width != 0 {
		sourceWidth = prop.width
	}

	sourceHeight := float64(dec.Bounds().Size().Y)
	if prop.height != 0 {
		sourceHeight = prop.height
	}




	//for i:= float64(0) ; i< prop.width;i++ {
	//	x :=  int(math.Round((i/prop.width)* sourceWidth))
	//	for j := float64(0);j<prop.height;j++ {
	//		y := int(math.Round((j/prop.height) * sourceHeight))
	//		col := dec.At(x,y)
	//		r,g,b,_ := col.RGBA()
	//
	//		//R ,G,B := float64(r)*0.299,float64(g)*0.587,float64(b)*0.114
	//		//tot := uint8((R+G+B)/257)
	//		//fmt.Println(R,G,B)
	//		//uintR ,uintG,uintb := uint8(r/257),uint8(g/257),uint8(b/257)
	//		//newImage.SetRGBA(int(i),int(j),color.RGBA{R: tot,G: tot,B: tot})
	//	}
	//}
	//fmt.Println("new image bounds", newImage.Bounds().Size())
	//draw.Draw(newImage,newImage.Bounds(),newImage,image.Pt(0,0),draw.Src)

	f, err := os.Create("3.jpg")
	if err != nil {
		log.Fatal(err)
	}


	jpeg.Encode(f,g,&jpeg.Options{Quality: 100})
	fmt.Println("done")
	return nil,nil
}
