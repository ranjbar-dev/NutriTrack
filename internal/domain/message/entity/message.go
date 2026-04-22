package entity

import (
	"time"

	"github.com/google/uuid"
)

// Message represents a chat message between a client and their nutritionist.
type Message struct {
	id             uuid.UUID
	senderID       uuid.UUID
	receiverID     uuid.UUID
	content        string
	attachmentPath *string
	attachmentType *string
	attachmentSize *int64
	attachmentName *string
	readAt         *time.Time
	createdAt      time.Time
}

// NewMessage creates a new Message for persistence (ID assigned by DB).
func NewMessage(senderID, receiverID uuid.UUID, content string) *Message {
	return &Message{
		senderID:   senderID,
		receiverID: receiverID,
		content:    content,
	}
}

// ReconstituteMessage rebuilds a Message from persisted storage data.
func ReconstituteMessage(id, senderID, receiverID uuid.UUID, content string, attachmentPath, attachmentType *string, attachmentSize *int64, attachmentName *string, readAt *time.Time, createdAt time.Time) *Message {
	return &Message{
		id:             id,
		senderID:       senderID,
		receiverID:     receiverID,
		content:        content,
		attachmentPath: attachmentPath,
		attachmentType: attachmentType,
		attachmentSize: attachmentSize,
		attachmentName: attachmentName,
		readAt:         readAt,
		createdAt:      createdAt,
	}
}

// Getters
func (m *Message) ID() uuid.UUID           { return m.id }
func (m *Message) SenderID() uuid.UUID     { return m.senderID }
func (m *Message) ReceiverID() uuid.UUID   { return m.receiverID }
func (m *Message) Content() string         { return m.content }
func (m *Message) AttachmentPath() *string { return m.attachmentPath }
func (m *Message) AttachmentType() *string { return m.attachmentType }
func (m *Message) AttachmentSize() *int64  { return m.attachmentSize }
func (m *Message) AttachmentName() *string { return m.attachmentName }
func (m *Message) ReadAt() *time.Time      { return m.readAt }
func (m *Message) CreatedAt() time.Time    { return m.createdAt }

// Setters
func (m *Message) SetAttachmentPath(v *string) { m.attachmentPath = v }
func (m *Message) SetAttachmentType(v *string) { m.attachmentType = v }
func (m *Message) SetAttachmentSize(v *int64)  { m.attachmentSize = v }
func (m *Message) SetAttachmentName(v *string) { m.attachmentName = v }

// HasAttachment returns true if the message includes a file attachment.
func (m *Message) HasAttachment() bool {
	return m.attachmentPath != nil
}

// MarkRead records when the message was read.
func (m *Message) MarkRead(at time.Time) {
	m.readAt = &at
}
