package imageconversion

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
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
		fmt.Printf("Get request failed to: %s \n", url)
		fmt.Println(err)
		return fmt.Errorf("Couldn't get the image from [the attachment url](%s)", url)
	}
	defer response.Body.Close()

	image, _, err := image.Decode(response.Body)
	if err != nil {
		// try decoding the image as a webp
		image, err = webp.Decode(response.Body, &decoder.Options{})
		if err != nil {
			// if it fails too, return the error
			fmt.Printf("Error decoding image from response body: %s \n", url)
			fmt.Println(err)
			return fmt.Errorf("I don't know how to handle the format of your image. (%s)", response.Header.Get("Content-Type"))
		}
	}

	imageResponse.Image = image
	// get the filename from the url
	urlStrings := strings.Split(response.Request.URL.Path, "/")
	imageResponse.Filename = urlStrings[len(urlStrings)-1]
	// only get the filename without the extension
	lastDot := strings.LastIndex(imageResponse.Filename, ".")
	if lastDot != -1 {
		// if there is an extension, remove it
		imageResponse.Filename = imageResponse.Filename[:lastDot]
	}

	return nil
}

func ConvertImage(format string, imageurl string) (ImageResponse, error) {
	err := getImageFromUrl(imageurl)
	if err != nil {
		return imageResponse, err
	}

	imageResponse.Filepath = fmt.Sprintf("%s/%s.%s", imageDirectory, imageResponse.Filename, format)
	imageResponse.Filename = fmt.Sprintf("%s.%s", imageResponse.Filename, format)
	buffer := new(bytes.Buffer)

	switch format {

	case "jpeg":
		err := jpeg.Encode(buffer, imageResponse.Image, &jpeg.Options{
			Quality: 90, // Quality factor (0:small..100:big)
		})
		if err != nil {
			fmt.Println("Error encoding the image to jpeg: ", err)
			return imageResponse, fmt.Errorf("I wasn't able to encode the image to jpeg...")
		}
		imageResponse.ContentType = "image/jpeg"
		imageResponse.Buffer = buffer
		return imageResponse, nil

	case "png":
		err := png.Encode(buffer, imageResponse.Image)
		if err != nil {
			fmt.Println("Error encoding the image to png: ", err)
			return imageResponse, fmt.Errorf("I wasn't able to encode the image to png...")
		}
		imageResponse.ContentType = "image/png"
		imageResponse.Buffer = buffer
		return imageResponse, nil

	case "webp":
		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 90)
		if err != nil {
			fmt.Println("Error creating the webp encoder options: ", err)
			return imageResponse, fmt.Errorf("I wasn't able to create the webp encoder options...")
		}
		err = webp.Encode(buffer, imageResponse.Image, options)
		if err != nil {
			fmt.Println("Error encoding the image to webp: ", err)
			return imageResponse, fmt.Errorf("I wasn't able to encode the image to webp...")
		}
		imageResponse.ContentType = "image/webp"
		imageResponse.Buffer = buffer
		return imageResponse, nil
	}

	return imageResponse, fmt.Errorf("The format %s is not supported", format)
}
