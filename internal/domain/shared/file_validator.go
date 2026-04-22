package shared

import "bytes"

// Supported image magic bytes
var (
	jpegMagic = []byte{0xFF, 0xD8, 0xFF}
	pngMagic  = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	webpMagic = []byte{0x52, 0x49, 0x46, 0x46} // RIFF....WEBP
)

// ImageInfo holds the detected image format metadata.
type ImageInfo struct {
	extension string // "jpg", "png", "webp"
	mimeType  string // "image/jpeg", "image/png", "image/webp"
}

// NewImageInfo constructs an ImageInfo value object.
func NewImageInfo(extension, mimeType string) *ImageInfo {
	return &ImageInfo{extension: extension, mimeType: mimeType}
}

// Extension returns the file extension (e.g. "jpg").
func (i *ImageInfo) Extension() string { return i.extension }

// MIMEType returns the MIME type (e.g. "image/jpeg").
func (i *ImageInfo) MIMEType() string { return i.mimeType }

// ValidateImageMagicBytes checks the first bytes of a file.
// Returns ImageInfo and nil error if valid; ErrInvalidFileType otherwise.
// header must be the first 12 bytes of the file (or all bytes if file is shorter).
func ValidateImageMagicBytes(header []byte) (*ImageInfo, error) {
	if bytes.HasPrefix(header, jpegMagic) {
		return NewImageInfo("jpg", "image/jpeg"), nil
	}
	if bytes.HasPrefix(header, pngMagic) {
		return NewImageInfo("png", "image/png"), nil
	}
	if bytes.HasPrefix(header, webpMagic) && len(header) >= 12 && string(header[8:12]) == "WEBP" {
		return NewImageInfo("webp", "image/webp"), nil
	}
	return nil, ErrInvalidFileType
}
