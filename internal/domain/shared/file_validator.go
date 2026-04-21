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
	Extension string // "jpg", "png", "webp"
	MIMEType  string // "image/jpeg", "image/png", "image/webp"
}

// ValidateImageMagicBytes checks the first bytes of a file.
// Returns ImageInfo and nil error if valid; ErrInvalidFileType otherwise.
// header must be the first 12 bytes of the file (or all bytes if file is shorter).
func ValidateImageMagicBytes(header []byte) (*ImageInfo, error) {
	if bytes.HasPrefix(header, jpegMagic) {
		return &ImageInfo{Extension: "jpg", MIMEType: "image/jpeg"}, nil
	}
	if bytes.HasPrefix(header, pngMagic) {
		return &ImageInfo{Extension: "png", MIMEType: "image/png"}, nil
	}
	if bytes.HasPrefix(header, webpMagic) && len(header) >= 12 && string(header[8:12]) == "WEBP" {
		return &ImageInfo{Extension: "webp", MIMEType: "image/webp"}, nil
	}
	return nil, ErrInvalidFileType
}
