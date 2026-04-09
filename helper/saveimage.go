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

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	path := filepath.Join(dir, filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		return "", err
	}

	return filename, nil
}

func SaveImages(c *gin.Context, formKey string, dir string) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	files := form.File[formKey]
	if len(files) == 0 {
		return []string{}, nil
	}

	// create folder if not exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	var paths []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		path := filepath.Join(dir, filename)

		if err := c.SaveUploadedFile(file, path); err != nil {
			return nil, err
		}
		paths = append(paths, filename)
	}

	return paths, nil
}

func DeleteFile(path string) {
	if path == "" {
		return
	}
	_ = os.Remove(path) // ignore error
}

func DeleteFiles(paths []string) {
	for _, p := range paths {
		_ = os.Remove(p)
	}
}
