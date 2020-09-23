package lib

import (
	_ "errors"
	"io"
	_ "log"
	"mime/multipart"
	"os"

)

func UploadImage(file *multipart.FileHeader, target string) error {
	// Open file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Prepare destination file
	dst, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// err = CompressImage(target, target)
	// if err != nil {
	// 	log.Println("Processing without compressing image:", err)
	// }

	return nil
}

// func CompressImage(input string, output string) error {
// 	options := bimg.Options{
// 		Quality:       90,
// 		Compression:   9,
// 		StripMetadata: true,
// 	}

// 	buffer, err := bimg.Read(input)
// 	if err != nil {
// 		return err
// 	}

// 	// Image type validation
// 	imageType := bimg.NewImage(buffer).Type()
// 	if (imageType != "png") && (imageType != "jpeg") {
// 		return errors.New("Unsupported type, must be png or jpeg")
// 	} else if imageType == "png" {
// 		// Check if it is animated PNG
// 		// Read acTL chunk, if it exist the image is animated
// 		// Reference: https://wiki.mozilla.org/APNG_Specification
// 		//log.Printf("%0x %0x %0x %0x", buffer[36], buffer[37], buffer[38], buffer[39])
// 		if (buffer[36] == '\x08') && (buffer[37] == '\x61') && (buffer[38] == '\x63') && (buffer[39] == '\x54') {
// 			return errors.New("Unsupported type, must be static png, not animated png")
// 		}
// 	}

// 	newImage, err := bimg.NewImage(buffer).Process(options)
// 	if err != nil {
// 		return err
// 	}

// 	err = bimg.Write(output, newImage)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }