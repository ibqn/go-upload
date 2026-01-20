package controllers

import (
	"go-upload/models"
	"go-upload/utils"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HandleGetFile(c *gin.Context) {

	fileId := c.Param("id")

	var upload models.Upload
	if err := utils.DB.Where("id = ?", fileId).First(&upload).Error; err != nil {
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
