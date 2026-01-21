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
	ErrImageProcessing   = errors.New("image processing failed")
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

const (
	DefaultQuality = 80
	MaxQuality     = 100
)

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

	opts := parseImageParams(c)

	// Debug logging
	fmt.Printf("Image request - ID: %s, Width: %d, Quality: %d, Format: %s\n", fileId, opts.Width, opts.Quality, opts.Format)

	var upload models.Upload
	if err := utils.DB.Where("id = ?", fileId).First(&upload).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(upload.FilePath))
	
	if !strings.HasPrefix(mimeType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is not an image"})
		return
	}

	optimizedImage, newMimeType, err := processImageWithOptions(upload.FilePath, opts, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
		return
	}

	c.Header("Content-Type", newMimeType)

	ext := getFileExtension(newMimeType)
	filename := fmt.Sprintf("image_%s.%s", fileId, ext)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))

	c.Data(http.StatusOK, newMimeType, optimizedImage)
}

func processImageWithOptions(filePath string, opts ImageOptions, contentType string) ([]byte, string, error) {
	// Debug logging
	fmt.Printf("processImageWithOptions called - filePath: %s, width: %d, quality: %d, format: %s, contentType: %s\n",
		filePath, opts.Width, opts.Quality, opts.Format, contentType)

	// Read the image file
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %v", err)
	}
	fmt.Printf("Image file read successfully, size: %d bytes\n", len(imageData))

	img := bimg.NewImage(imageData)
	if img == nil {
		return nil, "", fmt.Errorf("%w: failed to create bimg image", ErrImageProcessing)
	}

	size, err := img.Size()
	if err != nil {
		return nil, "", fmt.Errorf("%w: failed to get image size: %v", ErrImageProcessing, err)
	}
	fmt.Printf("Original image size: %dx%d\n", size.Width, size.Height)

	options := bimg.Options{
		Quality: opts.Quality,
	}

	if opts.Width > 0 {
		options.Width = opts.Width
		options.Enlarge = false
		options.Force = false
		fmt.Printf("Resizing to width: %d\n", opts.Width)
	}

	outputMimeType := contentType
	if opts.Format != "" && opts.Format != "original" {
		if imageType, supported := supportedFormats[opts.Format]; supported {
			options.Type = imageType
			outputMimeType = fmt.Sprintf("image/%s", opts.Format)
			if opts.Format == "jpg" {
				outputMimeType = "image/jpeg"
			}
			fmt.Printf("Converting format to: %s (mime: %s)\n", opts.Format, outputMimeType)
		} else {
			return nil, "", fmt.Errorf("%w: %s", ErrUnsupportedFormat, opts.Format)
		}
	} else {
		// Apply quality to original format
		if imageType, exists := mimeTypeMap[contentType]; exists {
			options.Type = imageType
			fmt.Printf("Applying quality %d to existing format: %s\n", opts.Quality, contentType)
		}
	}

	fmt.Printf("Processing with options: Quality=%d, Width=%d, Type=%v\n", options.Quality, options.Width, options.Type)

	processedImage, err := img.Process(options)
	if err != nil {
		return nil, "", fmt.Errorf("%w: bimg processing failed: %v", ErrImageProcessing, err)
	}

	if len(processedImage) == 0 {
		return nil, "", fmt.Errorf("%w: processed image is empty", ErrImageProcessing)
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
