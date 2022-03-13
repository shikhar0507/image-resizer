# Image-resizer
On the fly image Resizer with filters

### Supported Image types
1. jpeg/jpg
2. png
3. webp (**Note some of filters might not work on webp like negative**)

### How to use 
### Base url : `https://d3078njhubik3z.cloudfront.net/staging/imageResizer?`

### Required Params
1. Width `width=<anyWidth>`
2. Height `height=<anyHeight>`
3. Url `url=<imageUrl>`

## Example:

<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=800&height=800&url=https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=300&url=https://upload.wikimedia.org/wikipedia/commons/thumb/d/d3/0322_DNA_Nucleotides.jpg/370px-0322_DNA_Nucleotides.jpg>

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
<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=brightness&value=100>

4. ### Black & White
<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=blackAndWhite&value=100>

5. ### Blur
<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=blur&value=7>

6. ### Negative

<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=negative>

7. ### Positive

<https://d3078njhubik3z.cloudfront.net/staging/imageResizer?width=400&height=400&url=https://centralrecorder.com/wp-content/uploads/2021/07/claymore-season-2-release-810x456-1.jpg&filter=positive>
