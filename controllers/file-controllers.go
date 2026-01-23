package controllers

import (
	"go-upload/models"
	"go-upload/utils"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleGetFile(c *gin.Context) {
	fileId := c.Param("id")

	fileUUID, err := uuid.Parse(fileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file id"})
		return
	}

	var upload models.Upload
	if err := utils.DB.Where("id = ?", fileUUID).First(&upload).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(upload.FilePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	c.Header("Content-Type", mimeType)
	c.File(upload.FilePath)
}
