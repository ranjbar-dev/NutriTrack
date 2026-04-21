package db

import (
	"context"

	"github.com/google/uuid"
)

const createLabResult = `-- name: CreateLabResult :one
INSERT INTO lab_results (client_id, nutritionist_id, file_path, original_name, file_type, file_size, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, client_id, nutritionist_id, file_path, original_name, file_type, file_size, notes, created_at`

// CreateLabResultParams holds parameters for creating a lab result record.
type CreateLabResultParams struct {
	ClientID       uuid.UUID `db:"client_id"`
	NutritionistID uuid.UUID `db:"nutritionist_id"`
	FilePath       string    `db:"file_path"`
	OriginalName   string    `db:"original_name"`
	FileType       string    `db:"file_type"`
	FileSize       int64     `db:"file_size"`
	Notes          string    `db:"notes"`
}

// CreateLabResult inserts a new lab result record and returns the created row.
func (q *Queries) CreateLabResult(ctx context.Context, arg CreateLabResultParams) (LabResult, error) {
	row := q.db.QueryRow(ctx, createLabResult,
		arg.ClientID, arg.NutritionistID, arg.FilePath, arg.OriginalName, arg.FileType, arg.FileSize, arg.Notes,
	)
	var i LabResult
	err := row.Scan(
		&i.ID, &i.ClientID, &i.NutritionistID, &i.FilePath, &i.OriginalName, &i.FileType, &i.FileSize, &i.Notes, &i.CreatedAt,
	)
	return i, err
}

const getLabResultByID = `-- name: GetLabResultByID :one
SELECT id, client_id, nutritionist_id, file_path, original_name, file_type, file_size, notes, created_at
FROM lab_results WHERE id = $1`

// GetLabResultByID retrieves a lab result by its ID.
func (q *Queries) GetLabResultByID(ctx context.Context, id uuid.UUID) (LabResult, error) {
	row := q.db.QueryRow(ctx, getLabResultByID, id)
	var i LabResult
	err := row.Scan(
		&i.ID, &i.ClientID, &i.NutritionistID, &i.FilePath, &i.OriginalName, &i.FileType, &i.FileSize, &i.Notes, &i.CreatedAt,
	)
	return i, err
}

// ListLabResultsByClientIDParams holds parameters for paginated listing.
type ListLabResultsByClientIDParams struct {
	ClientID uuid.UUID `db:"client_id"`
	Limit    int32     `db:"limit"`
	Offset   int32     `db:"offset"`
}

const listLabResultsByClientID = `-- name: ListLabResultsByClientID :many
SELECT id, client_id, nutritionist_id, file_path, original_name, file_type, file_size, notes, created_at
FROM lab_results WHERE client_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

// ListLabResultsByClientID returns paginated lab results for a client.
func (q *Queries) ListLabResultsByClientID(ctx context.Context, arg ListLabResultsByClientIDParams) ([]LabResult, error) {
	rows, err := q.db.Query(ctx, listLabResultsByClientID, arg.ClientID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LabResult
	for rows.Next() {
		var i LabResult
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.NutritionistID, &i.FilePath, &i.OriginalName, &i.FileType, &i.FileSize, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const countLabResultsByClientID = `-- name: CountLabResultsByClientID :one
SELECT COUNT(*) FROM lab_results WHERE client_id = $1`

// CountLabResultsByClientID returns the total number of lab results for a client.
func (q *Queries) CountLabResultsByClientID(ctx context.Context, clientID uuid.UUID) (int64, error) {
	var count int64
	err := q.db.QueryRow(ctx, countLabResultsByClientID, clientID).Scan(&count)
	return count, err
}
