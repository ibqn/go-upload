package dto

// FileResponse represents file metadata
type FileResponse struct {
	FilePath string `json:"filePath"`
	MimeType string `json:"mimeType"`
}
