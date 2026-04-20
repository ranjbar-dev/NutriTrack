package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)
// CommunicationRepository defines operations for messaging and food requests.
type CommunicationRepository interface {
	// Messages
	CreateMessage(ctx context.Context, senderID, receiverID uuid.UUID, content, attachmentType, attachmentPath, attachmentName *string) (*dto.MessageResponse, error)
	ListMessages(ctx context.Context, userA, userB uuid.UUID, limit, offset int32) ([]dto.MessageResponse, error)
	ListMessagesSince(ctx context.Context, userA, userB uuid.UUID, since time.Time) ([]dto.MessageResponse, error)
	MarkMessagesRead(ctx context.Context, receiverID, senderID uuid.UUID) error
	CountUnreadMessages(ctx context.Context, receiverID uuid.UUID) (int32, error)
	GetMessageByID(ctx context.Context, id uuid.UUID) (*sqlc.Message, error)

	// Food Requests
	CreateFoodRequest(ctx context.Context, requestedBy uuid.UUID, foodName string, description *string) (*dto.FoodRequestResponse, error)
	ListFoodRequestsByClient(ctx context.Context, clientID uuid.UUID) ([]dto.FoodRequestResponse, error)
	ListPendingFoodRequestsForNutritionist(ctx context.Context, nutritionistID uuid.UUID) ([]dto.FoodRequestResponse, error)
	GetFoodRequestByID(ctx context.Context, id uuid.UUID) (*sqlc.FoodRequest, error)
	ApproveFoodRequest(ctx context.Context, id, reviewerID uuid.UUID) (*dto.FoodRequestResponse, error)
	RejectFoodRequest(ctx context.Context, id uuid.UUID, reason *string, reviewerID uuid.UUID) (*dto.FoodRequestResponse, error)
}

type communicationRepository struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

// NewCommunicationRepository creates a CommunicationRepository backed by pgxpool.
func NewCommunicationRepository(pool *pgxpool.Pool) CommunicationRepository {
	return &communicationRepository{pool: pool, q: sqlc.New(pool)}
}

func (r *communicationRepository) CreateMessage(ctx context.Context, senderID, receiverID uuid.UUID, content, attachmentType, attachmentPath, attachmentName *string) (*dto.MessageResponse, error) {
	msg, err := r.q.CreateMessage(ctx, sqlc.CreateMessageParams{
		SenderID:       pgtype.UUID{Bytes: senderID, Valid: true},
		ReceiverID:     pgtype.UUID{Bytes: receiverID, Valid: true},
		Content:        textOrNull(content),
		AttachmentType: textOrNull(attachmentType),
		AttachmentPath: textOrNull(attachmentPath),
		AttachmentName: textOrNull(attachmentName),
	})
	if err != nil {
		return nil, err
	}
	return messageToDTO(msg), nil
}

func (r *communicationRepository) ListMessages(ctx context.Context, userA, userB uuid.UUID, limit, offset int32) ([]dto.MessageResponse, error) {
	msgs, err := r.q.ListMessages(ctx, sqlc.ListMessagesParams{
		SenderID:   pgtype.UUID{Bytes: userA, Valid: true},
		ReceiverID: pgtype.UUID{Bytes: userB, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, err
	}
	result := make([]dto.MessageResponse, len(msgs))
	for i, m := range msgs {
		result[i] = *messageToDTO(m)
	}
	return result, nil
}

func (r *communicationRepository) ListMessagesSince(ctx context.Context, userA, userB uuid.UUID, since time.Time) ([]dto.MessageResponse, error) {
	msgs, err := r.q.ListMessagesSince(ctx, sqlc.ListMessagesSinceParams{
		SenderID:   pgtype.UUID{Bytes: userA, Valid: true},
		ReceiverID: pgtype.UUID{Bytes: userB, Valid: true},
		SentAt:     pgtype.Timestamptz{Time: since, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	result := make([]dto.MessageResponse, len(msgs))
	for i, m := range msgs {
		result[i] = *messageToDTO(m)
	}
	return result, nil
}

func (r *communicationRepository) MarkMessagesRead(ctx context.Context, receiverID, senderID uuid.UUID) error {
	return r.q.MarkMessagesRead(ctx, sqlc.MarkMessagesReadParams{
		ReceiverID: pgtype.UUID{Bytes: receiverID, Valid: true},
		SenderID:   pgtype.UUID{Bytes: senderID, Valid: true},
	})
}

func (r *communicationRepository) CountUnreadMessages(ctx context.Context, receiverID uuid.UUID) (int32, error) {
	return r.q.CountUnreadMessages(ctx, pgtype.UUID{Bytes: receiverID, Valid: true})
}

func (r *communicationRepository) GetMessageByID(ctx context.Context, id uuid.UUID) (*sqlc.Message, error) {
	msg, err := r.q.GetMessageByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *communicationRepository) CreateFoodRequest(ctx context.Context, requestedBy uuid.UUID, foodName string, description *string) (*dto.FoodRequestResponse, error) {
	fr, err := r.q.CreateFoodRequest(ctx, sqlc.CreateFoodRequestParams{
		FoodName:    foodName,
		Description: textOrNull(description),
		RequestedBy: pgtype.UUID{Bytes: requestedBy, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return foodRequestToDTO(fr, nil), nil
}

func (r *communicationRepository) ListFoodRequestsByClient(ctx context.Context, clientID uuid.UUID) ([]dto.FoodRequestResponse, error) {
	frs, err := r.q.ListFoodRequestsByClient(ctx, pgtype.UUID{Bytes: clientID, Valid: true})
	if err != nil {
		return nil, err
	}
	result := make([]dto.FoodRequestResponse, len(frs))
	for i, fr := range frs {
		result[i] = *foodRequestToDTO(fr, nil)
	}
	return result, nil
}

func (r *communicationRepository) ListPendingFoodRequestsForNutritionist(ctx context.Context, nutritionistID uuid.UUID) ([]dto.FoodRequestResponse, error) {
	frs, err := r.q.ListPendingFoodRequestsForNutritionist(ctx, pgtype.UUID{Bytes: nutritionistID, Valid: true})
	if err != nil {
		return nil, err
	}
	result := make([]dto.FoodRequestResponse, len(frs))
	for i, fr := range frs {
		result[i] = *foodRequestToDTO(fr, nil)
	}
	return result, nil
}

func (r *communicationRepository) GetFoodRequestByID(ctx context.Context, id uuid.UUID) (*sqlc.FoodRequest, error) {
	fr, err := r.q.GetFoodRequestByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &fr, nil
}

func (r *communicationRepository) ApproveFoodRequest(ctx context.Context, id, reviewerID uuid.UUID) (*dto.FoodRequestResponse, error) {
	fr, err := r.q.ApproveFoodRequest(ctx, sqlc.ApproveFoodRequestParams{
		ID:         pgtype.UUID{Bytes: id, Valid: true},
		ReviewedBy: pgtype.UUID{Bytes: reviewerID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return foodRequestToDTO(fr, nil), nil
}

func (r *communicationRepository) RejectFoodRequest(ctx context.Context, id uuid.UUID, reason *string, reviewerID uuid.UUID) (*dto.FoodRequestResponse, error) {
	fr, err := r.q.RejectFoodRequest(ctx, sqlc.RejectFoodRequestParams{
		ID:              pgtype.UUID{Bytes: id, Valid: true},
		RejectionReason: textOrNull(reason),
		ReviewedBy:      pgtype.UUID{Bytes: reviewerID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return foodRequestToDTO(fr, nil), nil
}

// ─── mapping helpers ─────────────────────────────────────────────────────────

func textOrNull(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func messageToDTO(m sqlc.Message) *dto.MessageResponse {
	resp := &dto.MessageResponse{
		ID:         uuid.UUID(m.ID.Bytes).String(),
		SenderID:   uuid.UUID(m.SenderID.Bytes).String(),
		ReceiverID: uuid.UUID(m.ReceiverID.Bytes).String(),
		SentAt:     m.SentAt.Time,
	}
	if m.Content.Valid {
		s := m.Content.String
		resp.Content = &s
	}
	if m.AttachmentType.Valid {
		s := m.AttachmentType.String
		resp.AttachmentType = &s
	}
	if m.AttachmentPath.Valid {
		s := m.AttachmentPath.String
		resp.AttachmentPath = &s
	}
	if m.AttachmentName.Valid {
		s := m.AttachmentName.String
		resp.AttachmentName = &s
	}
	if m.ReadAt.Valid {
		t := m.ReadAt.Time
		resp.ReadAt = &t
	}
	return resp
}

func foodRequestToDTO(fr sqlc.FoodRequest, clientName *string) *dto.FoodRequestResponse {
	resp := &dto.FoodRequestResponse{
		ID:          uuid.UUID(fr.ID.Bytes).String(),
		FoodName:    fr.FoodName,
		Status:      string(fr.Status),
		RequestedBy: uuid.UUID(fr.RequestedBy.Bytes).String(),
		ClientName:  clientName,
		CreatedAt:   fr.CreatedAt.Time,
		UpdatedAt:   fr.UpdatedAt.Time,
	}
	if fr.Description.Valid {
		s := fr.Description.String
		resp.Description = &s
	}
	if fr.RejectionReason.Valid {
		s := fr.RejectionReason.String
		resp.RejectionReason = &s
	}
	if fr.ReviewedBy.Valid {
		s := uuid.UUID(fr.ReviewedBy.Bytes).String()
		resp.ReviewedBy = &s
	}
	return resp
}
