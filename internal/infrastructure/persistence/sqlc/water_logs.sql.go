package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const upsertWaterLog = `-- name: UpsertWaterLog :one
INSERT INTO water_logs (client_id, local_id, logged_at, logged_date, amount_ml, notes)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING id, client_id, local_id, logged_at, logged_date, amount_ml, notes, created_at, (xmax::text::bigint = 0) AS inserted`

// UpsertWaterLogParams holds parameters for upserting a water log.
type UpsertWaterLogParams struct {
	ClientID   uuid.UUID `db:"client_id"`
	LocalID    string    `db:"local_id"`
	LoggedAt   time.Time `db:"logged_at"`
	LoggedDate time.Time `db:"logged_date"`
	AmountMl   int32     `db:"amount_ml"`
	Notes      string    `db:"notes"`
}

// UpsertWaterLog inserts or updates a water log and returns the row plus whether it was newly inserted.
func (q *Queries) UpsertWaterLog(ctx context.Context, arg UpsertWaterLogParams) (WaterLog, bool, error) {
	row := q.db.QueryRow(ctx, upsertWaterLog,
		arg.ClientID, arg.LocalID, arg.LoggedAt, arg.LoggedDate, arg.AmountMl, arg.Notes,
	)
	var i WaterLog
	var inserted bool
	err := row.Scan(
		&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate, &i.AmountMl, &i.Notes, &i.CreatedAt,
		&inserted,
	)
	return i, inserted, err
}

// ListWaterLogsByClientAndDateParams holds parameters for listing water logs by client and date.
type ListWaterLogsByClientAndDateParams struct {
	ClientID   uuid.UUID `db:"client_id"`
	LoggedDate time.Time `db:"logged_date"`
}

const listWaterLogsByClientAndDate = `-- name: ListWaterLogsByClientAndDate :many
SELECT id, client_id, local_id, logged_at, logged_date, amount_ml, notes, created_at
FROM water_logs
WHERE client_id = $1 AND logged_date = $2
ORDER BY logged_at ASC`

// ListWaterLogsByClientAndDate returns all water logs for a client on a given date.
func (q *Queries) ListWaterLogsByClientAndDate(ctx context.Context, arg ListWaterLogsByClientAndDateParams) ([]WaterLog, error) {
	rows, err := q.db.Query(ctx, listWaterLogsByClientAndDate, arg.ClientID, arg.LoggedDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WaterLog
	for rows.Next() {
		var i WaterLog
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate, &i.AmountMl, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ListWaterLogsByClientParams holds parameters for paginated water log listing.
type ListWaterLogsByClientParams struct {
	ClientID uuid.UUID `db:"client_id"`
	Limit    int32     `db:"limit"`
	Offset   int32     `db:"offset"`
}

const listWaterLogsByClient = `-- name: ListWaterLogsByClient :many
SELECT id, client_id, local_id, logged_at, logged_date, amount_ml, notes, created_at
FROM water_logs
WHERE client_id = $1
ORDER BY logged_date DESC, logged_at DESC
LIMIT $2 OFFSET $3`

// ListWaterLogsByClient returns paginated water logs for a client.
func (q *Queries) ListWaterLogsByClient(ctx context.Context, arg ListWaterLogsByClientParams) ([]WaterLog, error) {
	rows, err := q.db.Query(ctx, listWaterLogsByClient, arg.ClientID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WaterLog
	for rows.Next() {
		var i WaterLog
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate, &i.AmountMl, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const countWaterLogsByClient = `-- name: CountWaterLogsByClient :one
SELECT COUNT(*) FROM water_logs WHERE client_id = $1`

// CountWaterLogsByClient returns the total count of water logs for a client.
func (q *Queries) CountWaterLogsByClient(ctx context.Context, clientID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countWaterLogsByClient, clientID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
