//
// Custom JPEG compressor
// Learning golang - Eric Tatchell
//

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

// with the help of https://stackoverflow.com/questions/49594259/reading-image-in-go
func openImageFromPath(imagePath string) (image.Image, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Unable to open image from path: " + imagePath)
		return nil, err
	}
	defer imageFile.Close()
	image, _, err := image.Decode(imageFile)
	if err != nil {
		fmt.Println("Generic decoding failed -- trying PNG")
	} else {
		return image, err
	}
	imageFile.Seek(0, 0)

	loadedImage, err := png.Decode(imageFile)
	if err != nil {
		fmt.Println("Unable to decode the image. Try a different file?")
	}

	return loadedImage, err

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//	var jpegImage JPEGImage = JPEGImage{}
	var imagePath string = os.Args[1]
	image, err := openImageFromPath(imagePath)
	if err != nil {
		fmt.Println("Error opening/decoding the image file.")
		os.Exit(1)
	}
	var ycbcr []byte = GetByteArray(image)
	err = os.WriteFile("test.eric", ycbcr, 0644)
	check(err)
}
