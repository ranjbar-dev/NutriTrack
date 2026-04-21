package entity

import (
	"time"

	"github.com/google/uuid"
)

// Message represents a chat message between a client and their nutritionist.
type Message struct {
	ID             uuid.UUID
	SenderID       uuid.UUID
	ReceiverID     uuid.UUID
	Content        string
	AttachmentPath *string    // URL path or nil
	AttachmentType *string    // MIME type or nil
	AttachmentSize *int64     // bytes or nil
	AttachmentName *string    // original filename or nil
	ReadAt         *time.Time
	CreatedAt      time.Time
}

// HasAttachment returns true if the message includes a file attachment.
func (m *Message) HasAttachment() bool {
	return m.AttachmentPath != nil
}
