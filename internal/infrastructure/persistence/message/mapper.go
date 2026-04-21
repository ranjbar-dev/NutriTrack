package message

import (
	"github.com/jackc/pgx/v5"
	"github.com/ranjbar-dev/nutritrack/internal/domain/message/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

func toDomain(row db.Message) *entity.Message {
	return &entity.Message{
		ID:             row.ID,
		SenderID:       row.SenderID,
		ReceiverID:     row.ReceiverID,
		Content:        row.Content,
		AttachmentPath: row.AttachmentPath,
		AttachmentType: row.AttachmentType,
		AttachmentSize: row.AttachmentSize,
		AttachmentName: row.AttachmentName,
		ReadAt:         row.ReadAt,
		CreatedAt:      row.CreatedAt,
	}
}

func isNotFound(err error) bool {
	return err == pgx.ErrNoRows
}
