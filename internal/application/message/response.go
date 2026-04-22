package message

import (
	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/message/entity"
)

// MapMessageResponse converts a Message entity to a JSON-serialisable map.
// is_mine is true when the callerID matches the message sender.
// Defined here so that handlers in the interfaces layer do not need to import the entity package.
func MapMessageResponse(m *entity.Message, callerID uuid.UUID) map[string]any {
	result := map[string]any{
		"id":          m.ID(),
		"sender_id":   m.SenderID(),
		"receiver_id": m.ReceiverID(),
		"content":     m.Content(),
		"is_mine":     m.SenderID() == callerID,
		"read_at":     m.ReadAt(),
		"created_at":  m.CreatedAt(),
	}
	if m.HasAttachment() {
		result["attachment"] = map[string]any{
			"url":  m.AttachmentPath(),
			"type": m.AttachmentType(),
			"size": m.AttachmentSize(),
			"name": m.AttachmentName(),
		}
	} else {
		result["attachment"] = nil
	}
	return result
}
