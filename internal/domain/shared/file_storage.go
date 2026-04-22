package shared

import "io"

// FileStorage is the domain port for saving files.
// The LocalStorage adapter implements this interface in internal/infrastructure/storage/.
type FileStorage interface {
	SaveAvatar(src io.Reader, ext string) (string, error)
	SaveLabResult(src io.Reader, ext string) (string, error)
	SaveAttachment(src io.Reader, ext string) (string, error)
}
