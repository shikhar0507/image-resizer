# Image-resizer
On the fly image Resizer with filters. THis is a **FREE** application hosted on AWS Cloudfront. Use it as you like.

This uses nearest neighbour algorithm for faster image resizing.

### Supported Image types
1. jpeg/jpg
2. png
3. webp (**Note some of filters might not work on webp like negative**)

### How to use 
### Base url : `https://d3078njhubik3z.cloudfront.net/staging/imageResizer?`

### Required Params
1. Width `?width=<anyWidth>`
2. Height `?height=<anyHeight>`
3. Url `?url=<imageUrl>`
### Optional Params
1. filter `?filter=<filtername>`
2. value `?filter=<filtername>&value=<value>`

## Example:

<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=300&height=300&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg>


## Demo UI App
Refer to the README inside `panel` dir to know how to run the Demo UI App . That app is build with SVELTE 

## Filters

1. Grayscale
2. Sepia
3. Brightness  (**requires value**)
4. Black & White
5. Blur   (**requires value**)
6. Negative
7. Positive

### How to use

1. ### Grayscale

<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=grayscale>


2. ### Sepia
<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=sepia>

3. ### Brightness
Use a " - " value for `value` param to get a darker image. Ex : `...&value=-20`

<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=brightness&value=50>

4. ### Black & White
<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=blackAndWhite>

5. ### Blur
<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=blur&value=7>

6. ### Negative

<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=negative>

7. ### Positive

<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=positive>

## How to build locally 

1. Install GO
2. `go build main.go`
3. `./main`

### AWS Configuration

Please refer to the AWS Docs on AWS Sam, API Gateway,S3  etc..
