package controllers

import (
	"errors"
	"fmt"
	"go-upload/models"
	"go-upload/utils"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported image format")
	ErrInvalidParameters = errors.New("invalid parameters")
)

const (
	DefaultQuality = 80
	MaxQuality     = 100
)

var supportedFormats = map[string]bimg.ImageType{
	"webp": bimg.WEBP,
	"jpeg": bimg.JPEG,
	"jpg":  bimg.JPEG,
	"png":  bimg.PNG,
	"avif": bimg.AVIF,
}

var mimeTypeMap = map[string]bimg.ImageType{
	"image/jpeg": bimg.JPEG,
	"image/png":  bimg.PNG,
	"image/webp": bimg.WEBP,
}

type ImageOptions struct {
	Width   int
	Quality int
	Format  string
}

func parseImageParams(c *gin.Context) ImageOptions {
	opts := ImageOptions{
		Quality: DefaultQuality,
	}

	if w := c.Query("w"); w != "" {
		if width, err := strconv.Atoi(w); err == nil && width > 0 {
			opts.Width = width
		}
	}

	if q := c.Query("q"); q != "" {
		if quality, err := strconv.Atoi(q); err == nil && quality > 0 && quality <= MaxQuality {
			opts.Quality = quality
		}
	}

	opts.Format = c.Query("format")
	return opts
}

func HandleGetImage(c *gin.Context) {
	fileId := c.Param("id")

	width := c.Query("w")       // ?w=100 -> width
	quality := c.Query("q")     // ?q=80 -> quality
	format := c.Query("format") // ?format=webp -> format

	// Debug logging
	fmt.Printf("Image request - ID: %s, Width: %s, Quality: %s, Format: %s\n", fileId, width, quality, format)

	// With default values
	// width := c.DefaultQuery("w", "0")
	// quality := c.DefaultQuery("q", "80")
	// format := c.DefaultQuery("format", "original")

	var upload models.Upload
	if err := utils.DB.Where("id = ?", fileId).First(&upload).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(upload.FilePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if !strings.HasPrefix(mimeType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is not an image"})
		return
	}

	// Optimize image based on query parameters
	optimizedImage, newMimeType, err := optimizeImage(upload.FilePath, width, quality, format, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
		return
	}

	// Set content type header
	c.Header("Content-Type", newMimeType)

	// Set content disposition header with correct file extension
	ext := getFileExtension(newMimeType)
	filename := fmt.Sprintf("image_%s.%s", fileId, ext)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))

	c.Data(http.StatusOK, newMimeType, optimizedImage)
}

func optimizeImage(filePath string, widthStr string, qualityStr string, format string, contentType string) ([]byte, string, error) {
	// Debug logging
	fmt.Printf("optimizeImage called - filePath: %s, widthStr: %s, qualityStr: %s, format: %s, contentType: %s\n", filePath, widthStr, qualityStr, format, contentType)

	// Read the image file
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %v", err)
	}
	fmt.Printf("Image file read successfully, size: %d bytes\n", len(imageData))

	// Parse width parameter
	var width int
	if widthStr != "" {
		if w, err := strconv.Atoi(widthStr); err == nil && w > 0 {
			width = w
		}
	}

	// Parse quality parameter
	quality := 80 // default quality
	if qualityStr != "" {
		if q, err := strconv.Atoi(qualityStr); err == nil && q > 0 && q <= 100 {
			quality = q
		}
	}

	fmt.Printf("Parsed parameters - width: %d, quality: %d\n", width, quality)

	// Create bimg image
	img := bimg.NewImage(imageData)
	if img == nil {
		return nil, "", fmt.Errorf("failed to create bimg image")
	}

	// Get original image size
	size, err := img.Size()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get image size: %v", err)
	}
	fmt.Printf("Original image size: %dx%d\n", size.Width, size.Height)

	// Create bimg options
	options := bimg.Options{
		Quality: quality,
	}

	// Apply resize if width is provided
	if width > 0 {
		options.Width = width
		options.Enlarge = false // equivalent to withoutEnlargement: true
		options.Force = false   // equivalent to fit: "inside"
		fmt.Printf("Resizing to width: %d\n", width)
	}

	// Handle format conversion
	outputMimeType := contentType
	if format != "" {
		supportedFormats := map[string]bimg.ImageType{
			"webp": bimg.WEBP,
			"jpeg": bimg.JPEG,
			"jpg":  bimg.JPEG,
			"png":  bimg.PNG,
			"avif": bimg.AVIF,
		}

		if imageType, supported := supportedFormats[format]; supported {
			options.Type = imageType
			outputMimeType = fmt.Sprintf("image/%s", format)
			if format == "jpg" {
				outputMimeType = "image/jpeg"
			}
			fmt.Printf("Converting format to: %s (mime: %s)\n", format, outputMimeType)
		} else {
			return nil, "", fmt.Errorf("unsupported format: %s", format)
		}
	} else if qualityStr != "" && contentType != "" {
		// Apply quality to existing format
		formatMap := map[string]bimg.ImageType{
			"image/jpeg": bimg.JPEG,
			"image/png":  bimg.PNG,
			"image/webp": bimg.WEBP,
		}

		if imageType, exists := formatMap[contentType]; exists {
			options.Type = imageType
			fmt.Printf("Applying quality %d to existing format: %s\n", quality, contentType)
		}
	}

	fmt.Printf("Processing with options: Quality=%d, Width=%d, Type=%v\n", options.Quality, options.Width, options.Type)

	processedImage, err := img.Process(options)
	if err != nil {
		return nil, "", fmt.Errorf("bimg processing failed: %v", err)
	}

	if len(processedImage) == 0 {
		return nil, "", fmt.Errorf("processed image is empty")
	}

	fmt.Printf("Processing completed successfully, output size: %d bytes\n", len(processedImage))

	return processedImage, outputMimeType, nil
}

func getFileExtension(mimeType string) string {
	switch mimeType {
	case "image/webp":
		return "webp"
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/avif":
		return "avif"
	default:
		return "jpg"
	}
}
