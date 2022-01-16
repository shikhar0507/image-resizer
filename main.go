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
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Properties struct {
	url           *url.URL
	imageUrl      string
	width, height float64
	filter        string
	value         int
	//ext Extension
}

type responseHeaders struct {
	expires      string
	etag         string
	cacheControl string
	contentType  string
	lastModified time.Time
}

var property Properties
var ctx = context.Background()
var respHeader responseHeaders

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, evt events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	head := map[string]string{
		"Content-Type":  "image/jpeg",
		"Cache-Control": "max-age=63072000",
		"Accept-Ranges": "bytes",
	}
	heads := map[string]string{
		"Content-Type":  respHeader.contentType,
		"Cache-Control": respHeader.cacheControl,
		"ETag":          respHeader.etag,
		"Last-Modified": respHeader.lastModified.Local().String(),
	}

	req, err := handleRequest(evt)
	if err != nil {

		return events.APIGatewayProxyResponse{Headers: head, StatusCode: 500, Body: err.Error()}, err
	}
	fmt.Println("sending these headers", heads)
	return events.APIGatewayProxyResponse{Body: req, StatusCode: 200, Headers: heads, IsBase64Encoded: true}, nil
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

	ogUrl, _ := url.Parse(urlStr)
	property = Properties{url: ogUrl, width: float64(wd), height: float64(hg), filter: filter, imageUrl: urlStr, value: intVal}

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-2")})

	if err != nil {
		fmt.Fprintf(os.Stderr, "init error %s", err)
		return "", err
	}

	svc := s3.New(sess)
	input := &s3.HeadObjectInput{Bucket: aws.String("imageresizerbaseimages"), Key: aws.String(path.Base(urlStr))}
	_, err = svc.HeadObject(input)
	var resultImg image.Image

	if err == nil {

		img, imgErr := loadImageFromS3(svc, property)
		if imgErr != nil {
			fmt.Println("image found but error", imgErr)
			img, finalErr := loadImageFromUrl(property, sess)
			if finalErr != nil {
				return "", finalErr
			}
			resultImg = img
		} else {
			resultImg = img
		}
	} else {
		fmt.Println("image not found")
		img, imgErr := loadImageFromUrl(property, sess)
		if imgErr != nil {
			return "", imgErr
		}
		resultImg = img
	}

	if resultImg == nil {
		return "", errors.New("image not found")
	}

	buffer := new(bytes.Buffer)

	encodeErr := jpeg.Encode(buffer, resultImg, &jpeg.Options{Quality: 100})
	if encodeErr != nil {
		return "", encodeErr

	}
	encoded := base64.StdEncoding.EncodeToString([]byte(buffer.Bytes()))

	return encoded, nil
}

func loadImageFromS3(svc *s3.S3, prop Properties) (image.Image, error) {

	fmt.Println("getting image", path.Base(prop.imageUrl))
	ot, err := svc.GetObject(&s3.GetObjectInput{Bucket: aws.String("imageresizerbaseimages"), Key: aws.String(path.Base(prop.imageUrl))})

	if err != nil {
		fmt.Println("ot error", err.Error())
		return nil, err
	}
	respHeader.cacheControl = *ot.CacheControl
	respHeader.etag = *ot.ETag
	respHeader.contentType = *ot.ContentType
	respHeader.lastModified = *ot.LastModified

	fmt.Println(respHeader)
	defer ot.Body.Close()
	img, err := jpeg.Decode(ot.Body)
	if err != nil {
		fmt.Println("error decoding from s3", err.Error())
		return nil, err
	}
	return transformImage(img, prop)
}

func loadImageFromUrl(prop Properties, sess *session.Session) (image.Image, error) {

	respBody, err := getImage(prop.imageUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		return nil, err
	}

	saveToS3(respBody, path.Base(prop.url.String()), sess)
	decodedImage, err := jpeg.Decode(respBody)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		return nil, err
	}
	fmt.Println("decoded image")
	if decodedImage == nil {
		fmt.Fprintf(os.Stderr, "%s", "decoded image is nil")
		return nil, errors.New("decoded image is nill")
	}

	respHeader.cacheControl = "max-age:3600"
	respHeader.contentType = "image/jpeg"
	return transformImage(decodedImage, prop)
}

func saveToS3(data io.Reader, filename string, sess *session.Session) {
	log.Println("uploading image to s3")

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{Bucket: aws.String("imageresizerbaseimages"), Key: aws.String(filename), ContentType: aws.String("image/jpeg"), Body: data, CacheControl: aws.String("max-age=63072000")})

	if err != nil {
		fmt.Println("upload err", err.Error())
		return
	}
	log.Println("uploaded to", result.Location)
}

func getImage(imageUrl string) (io.Reader, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		return nil, err
	}

	return resp.Body, nil

}

func transformImage(decodedImage image.Image, prop Properties) (image.Image, error) {
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
