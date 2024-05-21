package imageconversion

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"
)

var (
	// The base path for the images to be saved to
	imageDirectory string        = "media/images"
	imageResponse  ImageResponse = ImageResponse{}
)

type ImageResponse struct {
	Filename    string
	ContentType string
	Filepath    string
	Image       image.Image
	Buffer      *bytes.Buffer
}

func getImageFromUrl(url string) error {
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		fmt.Printf("No response recieved from %s \n", url)
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()

	image, _, err := image.Decode(response.Body)
	if err != nil {
		fmt.Printf("Couldn't decode an image from the response body of %s \n", url)
		fmt.Println(err)
		return err
	}

	imageResponse.Image = image
	// get the filename from the url
	urlStrings := strings.Split(response.Request.URL.Path, "/")
	imageResponse.Filename = urlStrings[len(urlStrings)-1]
	// only get the filename without the extension
	lastDot := strings.LastIndex(imageResponse.Filename, ".")
	imageResponse.Filename = imageResponse.Filename[:lastDot]

	return nil
}

func ConvertImage(format string, imageurl string) (ImageResponse, error) {
	err := getImageFromUrl(imageurl)
	if err != nil {
		fmt.Printf("No image could be fetched from the url: %s\n", imageurl)
		return imageResponse, err
	}

	imageResponse.Filepath = fmt.Sprintf("%s/%s.%s", imageDirectory, imageResponse.Filename, format)
	imageResponse.Filename = fmt.Sprintf("%s.%s", imageResponse.Filename, format)

	switch format {

	case "jpeg":
		fmt.Println("Converting to jpeg")

		buffer := new(bytes.Buffer)
		err := jpeg.Encode(buffer, imageResponse.Image, &jpeg.Options{
			Quality: 90,
		})
		if err != nil {
			return imageResponse, err
		}
		imageResponse.ContentType = "image/jpeg"
		imageResponse.Buffer = buffer
    return imageResponse, nil

	case "png":
		fmt.Println("Converting to png")
		buffer := new(bytes.Buffer)
		err := png.Encode(buffer, imageResponse.Image)
		if err != nil {
			return imageResponse, err
		}
		imageResponse.ContentType = "image/png"
		imageResponse.Buffer = buffer
    return imageResponse, nil
	}

	return imageResponse, fmt.Errorf("The format %s is not supported", format)
}
