package shared

import "io"

// AttachmentStorage defines the storage port for message attachments.
// Implemented by internal/infrastructure/storage.LocalStorage.
type AttachmentStorage interface {
	SaveAttachment(src io.Reader, ext string) (string, error)
}

// LabResultStorage defines the storage port for lab result files.
// Implemented by internal/infrastructure/storage.LocalStorage.
type LabResultStorage interface {
	SaveLabResult(src io.Reader, ext string) (string, error)
}
