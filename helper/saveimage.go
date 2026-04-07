package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveImage(c *gin.Context, formKey string, dir string) (string, error) {
	file, err := c.FormFile(formKey)
	if err != nil {
		return "", err
	}

	// create folder if not exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	// get extension
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// generate unique name
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// full path
	path := filepath.Join(dir, filename)

	// save file
	if err := c.SaveUploadedFile(file, path); err != nil {
		return "", err
	}

	return path, nil
}
