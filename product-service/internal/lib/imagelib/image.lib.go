package imagelib

import (
	"io"
	"os"
)

func ReadDefaultImage() ([]byte, error) {
	filePath := "/assets/images/default_product_img.png"

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	imgByte, err := io.ReadAll(file)
	return imgByte, err
}
