package dto

// ImageParams represents image processing parameters
type ImageParams struct {
	Width   int    // Width in pixels (0 means no resize)
	Quality int    // Quality 1-100 (0 means default 80)
	Format  string // Format: jpeg, png, webp, avif (empty means original)
}
