package controllers

import (
	"fmt"
	"go-upload/models"
	"go-upload/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func CreateUpload(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := utils.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	fileSizeMB := file.Size / (1024 * 1024)
	if fileSizeMB > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
		return
	}

	folder := c.DefaultPostForm("folder", "")

	basePath := "file-storage"
	cleanFolder := filepath.Clean(folder)
	storagePath := filepath.Join(basePath, strconv.FormatUint(uint64(user.ID), 10), cleanFolder)

	if err := os.MkdirAll(storagePath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	fileName := file.Filename
	fullPath := filepath.Join(storagePath, fileName)

	if _, err := os.Stat(fullPath); err == nil {
		ext := filepath.Ext(fileName)
		nameWithoutExt := strings.TrimSuffix(fileName, ext)
		uniqueID := xid.New().String()
		fileName = fmt.Sprintf("%s_%s%s", nameWithoutExt, uniqueID, ext)
		fullPath = filepath.Join(storagePath, fileName)
	}

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	upload := models.Upload{
		UserID:   user.ID,
		FilePath: fullPath,
	}

	if err := utils.DB.Create(&upload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create upload record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"uploadId": upload.ID,
	})

}

func ListUploads(c *gin.Context) {}

func DeleteUpload(c *gin.Context) {}

func GetUpload(c *gin.Context) {}
