package imgconv

import (
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

func ConvertToBase64(file *multipart.FileHeader) (string, error) {
	mp, err := file.Open()
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(mp)
	if err != nil {
		return "", err
	}

	var b64Encoding string

	mimeType := http.DetectContentType(bytes)
	switch mimeType {
	case "image/bmp":
		b64Encoding += "data:image/bmp;base64,"
	case "image/webp":
		b64Encoding += "data:image/webp;base64,"
	case "image/png":
		b64Encoding += "data:image/png;base64,"
	case "image/jpeg":
		b64Encoding += "data:image/jpeg;base64,"
	case "image/gif":
		b64Encoding += "data:image/gif;base64,"
	default:
		return "", errors.New("not an image")
	}

	b64Encoding += base64.StdEncoding.EncodeToString(bytes)

	return b64Encoding, nil
}
