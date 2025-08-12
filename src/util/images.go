package util

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

const IMG_DIRECTORY = "/app/images/" // Should this be a configurable variable?

// / SaveImageB64 saves a b64 image with a MIME type onto permanent storage. /
// Returns (full_path, relative_path, error) / `full_path` is the full path on
// the device where the image was saved. / `relative_path` is the image path
// relative to IMG_DIRECTORY, including / file name and extension.
func SaveImageB64(base64Image string, filename string) (string, string, error) {
	// [1] Split payload

	parts := strings.Split(base64Image, ",")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid base64 image format")
	}

	mimeType := parts[0] // "data:image/png;base64"
	data := parts[1]     // Image (as base64)

	// [2] Decode image from B64

	imgBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", "", fmt.Errorf("could not decode base64 image, %v", err)
	}

	// [3] Determine file type

	var extension string
	if strings.Contains(mimeType, "image/png") {
		extension = ".png"
	} else if strings.Contains(mimeType, "image/jpeg") {
		extension = ".jpg"
	} else {
		return "", "", fmt.Errorf("unknown MIME type %s", mimeType)
	}

	// [4] Save image

	path_relative := filename + extension
	path := IMG_DIRECTORY + filename + extension
	err = os.WriteFile(path, imgBytes, 0644)
	if err != nil {
		return "", "", err
	}

	return path, path_relative, nil
}
