package helper

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

func ProtectImage(fileHeader *multipart.FileHeader) bool {

	// 1️⃣ Check extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return false
	}

	// 2️⃣ Open file
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	// 3️⃣ Read first 512 bytes (magic bytes)
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false
	}

	// 4️⃣ Detect real MIME type
	mimeType := http.DetectContentType(buffer)

	// 5️⃣ Allow only real images
	switch mimeType {
	case "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}
