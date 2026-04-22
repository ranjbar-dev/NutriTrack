package message

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/ranjbar-dev/nutritrack/internal/domain/message/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

func toDomain(row db.Message) *entity.Message {
	return entity.ReconstituteMessage(
		row.ID,
		row.SenderID,
		row.ReceiverID,
		row.Content,
		row.AttachmentPath,
		row.AttachmentType,
		row.AttachmentSize,
		row.AttachmentName,
		row.ReadAt,
		row.CreatedAt,
	)
}

func isNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
