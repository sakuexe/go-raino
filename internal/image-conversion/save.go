package imageconversion

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func createFile(filepath string) (*os.File, error) {
	// make sure that the directory exists
	err := os.MkdirAll(imageDirectory, os.ModePerm)
	if err != nil {
		fmt.Printf("Couldn't create the directory %s \n", imageDirectory)
		return nil, err
	}
	// save the image
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Couldn't create the file %s \n", filepath)
		return nil, err
	}

	return file, nil
}

func saveJpeg(image image.Image, filepath string) error {
	file, err := createFile(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := jpeg.Encode(file, image, nil); err != nil {
		fmt.Printf("Couldn't encode the image to the file %s \n", filepath)
		return err
	}

	return nil
}

func savePng(image image.Image, filepath string) error {
	file, err := createFile(filepath)
	if err != nil {
		return err
	}

	if err := png.Encode(file, image); err != nil {
		fmt.Printf("Couldn't encode the image to the file %s \n", filepath)
		return err
	}

	return nil
}
