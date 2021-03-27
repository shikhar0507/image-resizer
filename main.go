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
	//newImage := image.NewRGBA(image.Rect(0,0,int(prop.width),int(prop.height)))


	if err != nil {
		log.Fatal("error decoding jpeg",err)
	}

	//sourceWidth := float64(dec.Bounds().Size().X)
	//sourceHeight := float64(dec.Bounds().Size().Y)
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

	f, err := os.Create("i.jpg")
	if err != nil {
		log.Fatal(err)
	}

	//var opt Options
	//opt = Options{Quality: 75}
	imageGray := getImageGrayscale(dec,prop.width,prop.height)
	jpeg.Encode(f,blurImage(imageGray,prop.width,prop.height,9),&jpeg.Options{Quality: 100})
	fmt.Println("done")
	return nil,nil
}


func blurImage(dec image.Image,width float64,height float64,blurRadius int) image.Image {
	newImage := image.NewGray(image.Rect(0,0,int(width),int(height)))
	for x:= 0; x< int(width);x++ {
		for y:=0; y < int(height);y++ {
			//pos := dec.At(x,y)
			//blur radius loop
			totalValue := 0
			blurOff := (blurRadius-1)/2
			for i:= -blurOff;i <= blurOff;i++ {
				for j := -blurOff;j <= blurOff;j++ {
						pxA,_,_,_ := dec.At(x+i,y+j).RGBA()
						totalValue  = totalValue + (int(pxA) * 1)
				}
			}
			//fmt.Println(uint8((totalValue/(blurRadius)*2)/257))
			newImage.SetGray(x,y,color.Gray{Y: uint8((totalValue/257)/(blurRadius*blurRadius))})
		}
	}
	return newImage
}


func getImageGrayscale(dec image.Image,width float64,height float64) image.Image {
	newImage := image.NewGray(image.Rect(0,0,int(width),int(height)))
	sourceWidth := float64(dec.Bounds().Size().X)
	sourceHeight := float64(dec.Bounds().Size().Y)
	for i:= float64(0) ; i< width;i++ {
		x :=  int(math.Round((i/width)* sourceWidth))
		for j := float64(0);j<height;j++ {
			y := int(math.Round((j/height) * sourceHeight))
			col := dec.At(x,y)
			r,g,b,_ := col.RGBA()

			R ,G,B := float64(r)*0.299,float64(g)*0.587,float64(b)*0.114
			tot := uint8((R+G+B)/257)
			//fmt.Println(R,G,B)
			//uintR ,uintG,uintb := uint8(r/257),uint8(g/257),uint8(b/257)
			newImage.SetGray(int(i),int(j),color.Gray{Y:tot})
		}
	}
	return newImage
}
