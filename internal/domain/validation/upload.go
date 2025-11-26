package validation

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

const (
	// MaxImagesPerUpload is the maximum number of images allowed per upload
	MaxImagesPerUpload = 4

	// MaxImageSize is the maximum file size in bytes (10 MB)
	MaxImageSize = 10 * 1024 * 1024 // 10 MB

	// MinImageSize is the minimum file size in bytes (1 KB)
	MinImageSize = 1024 // 1 KB
)

// AllowedImageMimeTypes are the allowed MIME types for image uploads
var AllowedImageMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// AllowedImageExtensions are the allowed file extensions for image uploads
var AllowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// ValidateImageFile validates a single image file
func ValidateImageFile(file *multipart.FileHeader) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}

	// Check file size
	if file.Size > MaxImageSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxImageSize)
	}

	if file.Size < MinImageSize {
		return fmt.Errorf("file size is below minimum allowed size of %d bytes", MinImageSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedImageExtensions[ext] {
		return fmt.Errorf("file extension '%s' is not allowed. Allowed extensions: jpg, jpeg, png, gif, webp", ext)
	}

	// Check MIME type from content type header
	contentType := file.Header.Get("Content-Type")
	if contentType != "" && !AllowedImageMimeTypes[contentType] {
		return fmt.Errorf("file MIME type '%s' is not allowed. Allowed types: image/jpeg, image/png, image/gif, image/webp", contentType)
	}

	return nil
}

// ValidateImageFiles validates multiple image files
func ValidateImageFiles(files []*multipart.FileHeader, maxFiles int) error {
	if len(files) == 0 {
		return nil // No files to validate
	}

	if len(files) > maxFiles {
		return fmt.Errorf("too many files. Maximum %d files allowed, got %d", maxFiles, len(files))
	}

	for i, file := range files {
		if err := ValidateImageFile(file); err != nil {
			return fmt.Errorf("file %d validation failed: %w", i+1, err)
		}
	}

	return nil
}

// HasValidFilename checks if the file header has a non-empty filename
func HasValidFilename(file *multipart.FileHeader) bool {
	return file != nil && file.Filename != ""
}

// FilterValidFiles filters out files with empty filenames
func FilterValidFiles(files []*multipart.FileHeader) []*multipart.FileHeader {
	validFiles := make([]*multipart.FileHeader, 0, len(files))
	for _, file := range files {
		if HasValidFilename(file) {
			validFiles = append(validFiles, file)
		}
	}
	return validFiles
}
