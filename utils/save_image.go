package utils

import (
	"errors"

	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
)

var (
	ErrBucket       = errors.New("invalid bucket")
	ErrSize         = errors.New("invalid size")
	ErrInvalidImage = errors.New("invalid image")
)

func SaveImageToDisk(fileNameBase string, data []byte) (string, error) {
	_, fm, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	fileName := "./assets/property_images/" + fileNameBase + "." + fm
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
