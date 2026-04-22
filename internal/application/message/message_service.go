package message

import (
	"bytes"
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/message/entity"
	msgRepo "github.com/ranjbar-dev/nutritrack/internal/domain/message/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
)

const (
	maxAttachmentSizeImage = 5 * 1024 * 1024  // 5 MB for JPEG/PNG
	maxAttachmentSizePDF   = 10 * 1024 * 1024 // 10 MB for PDF
)

// allowedAttachmentMIME maps MIME to file extension.
var allowedAttachmentMIME = map[string]string{
	"application/pdf": "pdf",
	"image/jpeg":      "jpg",
	"image/png":       "png",
}

// MessageService provides message management business logic.
type MessageService struct {
	msgRepo  msgRepo.MessageRepository
	userRepo userRepo.UserRepository
	storage  shared.AttachmentStorage
}

// NewMessageService creates a new MessageService.
func NewMessageService(
	msgRepo msgRepo.MessageRepository,
	userRepo userRepo.UserRepository,
	storage shared.AttachmentStorage,
) *MessageService {
	return &MessageService{msgRepo: msgRepo, userRepo: userRepo, storage: storage}
}

// hasMagic checks if data starts with the given magic bytes.
func hasMagic(data, magic []byte) bool {
	if len(data) < len(magic) {
		return false
	}
	return bytes.Equal(data[:len(magic)], magic)
}

// detectAttachmentMIME reads magic bytes and returns MIME type + reconstituted reader.
func detectAttachmentMIME(r io.Reader) (string, io.Reader, error) {
	buf := make([]byte, 8)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", nil, shared.ErrInvalidFileType
	}
	buf = buf[:n]

	var mime string
	switch {
	case hasMagic(buf, []byte{0x25, 0x50, 0x44, 0x46}):
		mime = "application/pdf"
	case hasMagic(buf, []byte{0xFF, 0xD8, 0xFF}):
		mime = "image/jpeg"
	case hasMagic(buf, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}):
		mime = "image/png"
	default:
		return "", nil, shared.ErrInvalidFileType
	}

	full := io.MultiReader(bytes.NewReader(buf), r)
	return mime, full, nil
}

// saveAttachment validates, detects MIME, and saves an attachment file.
func (s *MessageService) saveAttachment(src io.Reader, originalName string, fileSize int64) (*string, *string, *int64, *string, error) {
	mimeType, fullReader, err := detectAttachmentMIME(src)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	ext, ok := allowedAttachmentMIME[mimeType]
	if !ok {
		return nil, nil, nil, nil, shared.ErrInvalidFileType
	}

	maxSize := int64(maxAttachmentSizePDF)
	if mimeType == "image/jpeg" || mimeType == "image/png" {
		maxSize = maxAttachmentSizeImage
	}
	if fileSize > maxSize {
		return nil, nil, nil, nil, shared.ErrFileTooLarge
	}

	urlPath, err := s.storage.SaveAttachment(fullReader, ext)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return &urlPath, &mimeType, &fileSize, &originalName, nil
}

// SendAsClient sends a message from a client to their nutritionist.
func (s *MessageService) SendAsClient(
	ctx context.Context,
	clientID uuid.UUID,
	content string,
	attachment io.Reader,
	attachmentName string,
	attachmentSize int64,
) (*entity.Message, error) {
	client, err := s.userRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, shared.ErrUserNotFound
	}
	if client.GetNutritionistID() == nil {
		return nil, shared.ErrInternal
	}

	msg := entity.NewMessage(clientID, *client.GetNutritionistID(), content)

	if attachment != nil {
		path, mimeType, size, name, saveErr := s.saveAttachment(attachment, attachmentName, attachmentSize)
		if saveErr != nil {
			return nil, saveErr
		}
		msg.SetAttachmentPath(path)
		msg.SetAttachmentType(mimeType)
		msg.SetAttachmentSize(size)
		msg.SetAttachmentName(name)
	}

	if msg.Content() == "" && msg.AttachmentPath() == nil {
		return nil, shared.ErrValidation
	}

	return s.msgRepo.Create(ctx, msg)
}

// SendAsNutritionist sends a message from a nutritionist to one of their clients.
func (s *MessageService) SendAsNutritionist(
	ctx context.Context,
	nutritionistID uuid.UUID,
	clientID uuid.UUID,
	content string,
	attachment io.Reader,
	attachmentName string,
	attachmentSize int64,
) (*entity.Message, error) {
	client, err := s.userRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, shared.ErrUserNotFound
	}
	if !client.BelongsTo(nutritionistID) {
		return nil, shared.ErrForbidden
	}

	msg := entity.NewMessage(nutritionistID, clientID, content)

	if attachment != nil {
		path, mimeType, size, name, saveErr := s.saveAttachment(attachment, attachmentName, attachmentSize)
		if saveErr != nil {
			return nil, saveErr
		}
		msg.SetAttachmentPath(path)
		msg.SetAttachmentType(mimeType)
		msg.SetAttachmentSize(size)
		msg.SetAttachmentName(name)
	}

	if msg.Content() == "" && msg.AttachmentPath() == nil {
		return nil, shared.ErrValidation
	}

	return s.msgRepo.Create(ctx, msg)
}

// GetClientConversation retrieves the conversation for a client with their nutritionist.
// Also marks the other party's messages as read.
func (s *MessageService) GetClientConversation(
	ctx context.Context,
	clientID uuid.UUID,
	limit, offset int32,
) ([]*entity.Message, int64, error) {
	client, err := s.userRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, 0, err
	}
	if client == nil || client.GetNutritionistID() == nil {
		return nil, 0, shared.ErrUserNotFound
	}
	nutritionistID := *client.GetNutritionistID()

	msgs, total, err := s.msgRepo.ListConversation(ctx, clientID, nutritionistID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Mark nutritionist's messages as read by client
	_ = s.msgRepo.MarkRead(ctx, clientID, nutritionistID)

	return msgs, total, nil
}

// GetNutritionistConversation retrieves the conversation between a nutritionist and one of their clients.
// Also marks the client's messages as read.
func (s *MessageService) GetNutritionistConversation(
	ctx context.Context,
	nutritionistID uuid.UUID,
	clientID uuid.UUID,
	limit, offset int32,
) ([]*entity.Message, int64, error) {
	client, err := s.userRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, 0, err
	}
	if client == nil {
		return nil, 0, shared.ErrUserNotFound
	}
	if !client.BelongsTo(nutritionistID) {
		return nil, 0, shared.ErrForbidden
	}

	msgs, total, err := s.msgRepo.ListConversation(ctx, nutritionistID, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Mark client's messages as read by nutritionist
	_ = s.msgRepo.MarkRead(ctx, nutritionistID, clientID)

	return msgs, total, nil
}

// GetUnreadCount returns the number of unread messages for the given user.
func (s *MessageService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.msgRepo.CountUnread(ctx, userID)
}
