package validation

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
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
		return exception.FileValidationError{Message: "file is nil"}
	}

	// Check file size
	if file.Size > MaxImageSize {
		return exception.FileValidationError{
			Message: fmt.Sprintf("file size exceeds maximum allowed size of %d bytes (%.2f MB)", MaxImageSize, float64(MaxImageSize)/(1024*1024)),
		}
	}

	if file.Size < MinImageSize {
		return exception.FileValidationError{
			Message: fmt.Sprintf("file size is below minimum allowed size of %d bytes", MinImageSize),
		}
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedImageExtensions[ext] {
		return exception.FileValidationError{
			Message: fmt.Sprintf("file extension '%s' is not allowed. Allowed extensions: jpg, jpeg, png, gif, webp", ext),
		}
	}

	// Check MIME type from content type header
	contentType := file.Header.Get("Content-Type")
	if contentType != "" && !AllowedImageMimeTypes[contentType] {
		return exception.FileValidationError{
			Message: fmt.Sprintf("file MIME type '%s' is not allowed. Allowed types: image/jpeg, image/png, image/gif, image/webp", contentType),
		}
	}

	return nil
}

// ValidateImageFiles validates multiple image files
func ValidateImageFiles(files []*multipart.FileHeader, maxFiles int) error {
	if len(files) == 0 {
		return nil // No files to validate
	}

	if len(files) > maxFiles {
		return exception.FileValidationError{
			Message: fmt.Sprintf("too many files. Maximum %d files allowed, got %d", maxFiles, len(files)),
		}
	}

	for i, file := range files {
		if err := ValidateImageFile(file); err != nil {
			// If it's already a FileValidationError, wrap it with file number
			if fileErr, ok := err.(exception.FileValidationError); ok {
				return exception.FileValidationError{
					Message: fmt.Sprintf("file %d validation failed: %s", i+1, fileErr.Message),
				}
			}
			return exception.FileValidationError{
				Message: fmt.Sprintf("file %d validation failed: %s", i+1, err.Error()),
			}
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
