package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const upsertMedicationLog = `-- name: UpsertMedicationLog :one
INSERT INTO medication_logs (client_id, local_id, logged_at, logged_date, medication_id, medication_name, dosage, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING id, client_id, local_id, logged_at, logged_date, medication_id, medication_name, dosage, notes, created_at, (xmax::text::bigint = 0) AS inserted`

// UpsertMedicationLogParams holds parameters for upserting a medication log.
type UpsertMedicationLogParams struct {
	ClientID       uuid.UUID  `db:"client_id"`
	LocalID        string     `db:"local_id"`
	LoggedAt       time.Time  `db:"logged_at"`
	LoggedDate     time.Time  `db:"logged_date"`
	MedicationID   *uuid.UUID `db:"medication_id"`
	MedicationName string     `db:"medication_name"`
	Dosage         string     `db:"dosage"`
	Notes          string     `db:"notes"`
}

// UpsertMedicationLog inserts or updates a medication log and returns the row plus whether it was newly inserted.
func (q *Queries) UpsertMedicationLog(ctx context.Context, arg UpsertMedicationLogParams) (MedicationLog, bool, error) {
	row := q.db.QueryRow(ctx, upsertMedicationLog,
		arg.ClientID, arg.LocalID, arg.LoggedAt, arg.LoggedDate,
		arg.MedicationID, arg.MedicationName, arg.Dosage, arg.Notes,
	)
	var i MedicationLog
	var inserted bool
	err := row.Scan(
		&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate,
		&i.MedicationID, &i.MedicationName, &i.Dosage, &i.Notes, &i.CreatedAt,
		&inserted,
	)
	return i, inserted, err
}

// ListMedicationLogsByClientAndDateParams holds parameters for listing medication logs by client and date.
type ListMedicationLogsByClientAndDateParams struct {
	ClientID   uuid.UUID `db:"client_id"`
	LoggedDate time.Time `db:"logged_date"`
}

const listMedicationLogsByClientAndDate = `-- name: ListMedicationLogsByClientAndDate :many
SELECT id, client_id, local_id, logged_at, logged_date, medication_id, medication_name, dosage, notes, created_at
FROM medication_logs
WHERE client_id = $1 AND logged_date = $2
ORDER BY logged_at ASC`

// ListMedicationLogsByClientAndDate returns all medication logs for a client on a given date.
func (q *Queries) ListMedicationLogsByClientAndDate(ctx context.Context, arg ListMedicationLogsByClientAndDateParams) ([]MedicationLog, error) {
	rows, err := q.db.Query(ctx, listMedicationLogsByClientAndDate, arg.ClientID, arg.LoggedDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MedicationLog
	for rows.Next() {
		var i MedicationLog
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate,
			&i.MedicationID, &i.MedicationName, &i.Dosage, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ListMedicationLogsByClientParams holds parameters for paginated medication log listing.
type ListMedicationLogsByClientParams struct {
	ClientID uuid.UUID `db:"client_id"`
	Limit    int32     `db:"limit"`
	Offset   int32     `db:"offset"`
}

const listMedicationLogsByClient = `-- name: ListMedicationLogsByClient :many
SELECT id, client_id, local_id, logged_at, logged_date, medication_id, medication_name, dosage, notes, created_at
FROM medication_logs
WHERE client_id = $1
ORDER BY logged_date DESC, logged_at DESC
LIMIT $2 OFFSET $3`

// ListMedicationLogsByClient returns paginated medication logs for a client.
func (q *Queries) ListMedicationLogsByClient(ctx context.Context, arg ListMedicationLogsByClientParams) ([]MedicationLog, error) {
	rows, err := q.db.Query(ctx, listMedicationLogsByClient, arg.ClientID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MedicationLog
	for rows.Next() {
		var i MedicationLog
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate,
			&i.MedicationID, &i.MedicationName, &i.Dosage, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const countMedicationLogsByClient = `-- name: CountMedicationLogsByClient :one
SELECT COUNT(*) FROM medication_logs WHERE client_id = $1`

// CountMedicationLogsByClient returns the total count of medication logs for a client.
func (q *Queries) CountMedicationLogsByClient(ctx context.Context, clientID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countMedicationLogsByClient, clientID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
