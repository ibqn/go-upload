package controllers

import (
	"fmt"
	"go-upload/models"
	"go-upload/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/xid"
)

func CreateUpload(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
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
	storagePath := filepath.Join(basePath, userUUID.String(), cleanFolder)

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
		UserID:   userUUID,
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

func ListUploads(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var uploads []models.Upload
	if err := utils.DB.Where("user_id = ?", userUUID).Find(&uploads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve uploads"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Uploads retrieved successfully",
		"uploads": uploads})
}

func DeleteUpload(c *gin.Context) {
	uploadId := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	uploadUUID, err := uuid.Parse(uploadId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload id"})
		return
	}

	var upload models.Upload
	if err := utils.DB.Where("id = ?", uploadUUID).First(&upload).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload not found"})
		return
	}

	if upload.UserID != userUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this upload"})
		return
	}

	if err := os.Remove(upload.FilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from storage"})
		return
	}

	if err := utils.DB.Delete(&upload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete upload record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload deleted successfully"})
}

func GetUpload(c *gin.Context) {
	uploadId := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	uploadUUID, err := uuid.Parse(uploadId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload id"})
		return
	}

	var upload models.Upload
	if err := utils.DB.Where("id = ?", uploadUUID).First(&upload).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload not found"})
		return
	}

	if upload.UserID != userUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this upload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload retrieved successfully",
		"upload":  upload,
	})
}
