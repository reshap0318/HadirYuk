package face

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg" // register JPEG decoder
	_ "image/png"  // register PNG decoder
	"os"
	"path/filepath"
	"strings"

	"github.com/reshap0318/hadirYuk/internal/helpers"
)

var allowedFaceImageExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

const maxFacePhotoSize = 2 * 1024 * 1024 // 2MB

// Base64ToTempFile decodes a base64-encoded image and writes it to a temporary file.
// Returns the path to the temporary file.
func Base64ToTempFile(base64Data string) (string, error) {
	if base64Data == "" {
		return "", fmt.Errorf("base64 data is empty")
	}

	// Strip data URI prefix if present (e.g., "data:image/jpeg;base64,")
	if idx := strings.Index(base64Data, ","); idx != -1 {
		base64Data = base64Data[idx+1:]
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %w", err)
	}

	if len(data) > maxFacePhotoSize {
		return "", fmt.Errorf("image size exceeds maximum allowed size of 2MB")
	}

	// Detect image format from header to determine extension
	_, format, err := image.DecodeConfig(strings.NewReader(string(data)))
	if err != nil {
		return "", fmt.Errorf("invalid image format: %w", err)
	}

	ext := "." + strings.ToLower(format)
	if !allowedFaceImageExts[ext] {
		return "", fmt.Errorf("image format %s is not allowed, use JPG/PNG/WebP", ext)
	}

	tmpDir := "storage/tmp"
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	randomStr, err := helpers.GenerateRandomString(8)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	fileName := fmt.Sprintf("face_%s%s", randomStr, ext)
	filePath := filepath.Join(tmpDir, fileName)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	return filePath, nil
}

// ValidateFaceImageFormat checks if the image file exists, has a valid format, and size.
func ValidateFaceImageFormat(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path is empty")
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return fmt.Errorf("failed to stat file: %w", err)
	}

	if info.Size() > maxFacePhotoSize {
		return fmt.Errorf("file size exceeds maximum allowed size of 2MB")
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	if !allowedFaceImageExts[ext] {
		return fmt.Errorf("file format %s is not allowed, use JPG/PNG/WebP", ext)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	_, _, err = image.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("invalid image file: %w", err)
	}

	return nil
}

// CleanupTempFile removes a temporary file if it exists.
func CleanupTempFile(filePath string) {
	if filePath != "" {
		_ = os.Remove(filePath)
	}
}
