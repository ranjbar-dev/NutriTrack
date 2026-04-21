package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const upsertSleepLog = `-- name: UpsertSleepLog :one
INSERT INTO sleep_logs (client_id, local_id, logged_date, sleep_start, sleep_end, duration_minutes, quality, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING id, client_id, local_id, logged_date, sleep_start, sleep_end, duration_minutes, quality, notes, created_at, (xmax::text::bigint = 0) AS inserted`

// UpsertSleepLogParams holds parameters for upserting a sleep log.
type UpsertSleepLogParams struct {
	ClientID        uuid.UUID `db:"client_id"`
	LocalID         string    `db:"local_id"`
	LoggedDate      time.Time `db:"logged_date"`
	SleepStart      time.Time `db:"sleep_start"`
	SleepEnd        time.Time `db:"sleep_end"`
	DurationMinutes int32     `db:"duration_minutes"`
	Quality         int32     `db:"quality"`
	Notes           string    `db:"notes"`
}

// UpsertSleepLog inserts or updates a sleep log and returns the row plus whether it was newly inserted.
func (q *Queries) UpsertSleepLog(ctx context.Context, arg UpsertSleepLogParams) (SleepLog, bool, error) {
	row := q.db.QueryRow(ctx, upsertSleepLog,
		arg.ClientID, arg.LocalID, arg.LoggedDate, arg.SleepStart, arg.SleepEnd,
		arg.DurationMinutes, arg.Quality, arg.Notes,
	)
	var i SleepLog
	var inserted bool
	err := row.Scan(
		&i.ID, &i.ClientID, &i.LocalID, &i.LoggedDate, &i.SleepStart, &i.SleepEnd,
		&i.DurationMinutes, &i.Quality, &i.Notes, &i.CreatedAt,
		&inserted,
	)
	return i, inserted, err
}

// ListSleepLogsByClientAndDateParams holds parameters for listing sleep logs by client and date.
type ListSleepLogsByClientAndDateParams struct {
	ClientID   uuid.UUID `db:"client_id"`
	LoggedDate time.Time `db:"logged_date"`
}

const listSleepLogsByClientAndDate = `-- name: ListSleepLogsByClientAndDate :many
SELECT id, client_id, local_id, logged_date, sleep_start, sleep_end, duration_minutes, quality, notes, created_at
FROM sleep_logs
WHERE client_id = $1 AND logged_date = $2
ORDER BY sleep_start ASC`

// ListSleepLogsByClientAndDate returns all sleep logs for a client on a given date.
func (q *Queries) ListSleepLogsByClientAndDate(ctx context.Context, arg ListSleepLogsByClientAndDateParams) ([]SleepLog, error) {
	rows, err := q.db.Query(ctx, listSleepLogsByClientAndDate, arg.ClientID, arg.LoggedDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SleepLog
	for rows.Next() {
		var i SleepLog
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.LoggedDate, &i.SleepStart, &i.SleepEnd,
			&i.DurationMinutes, &i.Quality, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ListSleepLogsByClientParams holds parameters for paginated sleep log listing.
type ListSleepLogsByClientParams struct {
	ClientID uuid.UUID `db:"client_id"`
	Limit    int32     `db:"limit"`
	Offset   int32     `db:"offset"`
}

const listSleepLogsByClient = `-- name: ListSleepLogsByClient :many
SELECT id, client_id, local_id, logged_date, sleep_start, sleep_end, duration_minutes, quality, notes, created_at
FROM sleep_logs
WHERE client_id = $1
ORDER BY logged_date DESC, sleep_start DESC
LIMIT $2 OFFSET $3`

// ListSleepLogsByClient returns paginated sleep logs for a client.
func (q *Queries) ListSleepLogsByClient(ctx context.Context, arg ListSleepLogsByClientParams) ([]SleepLog, error) {
	rows, err := q.db.Query(ctx, listSleepLogsByClient, arg.ClientID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SleepLog
	for rows.Next() {
		var i SleepLog
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.LoggedDate, &i.SleepStart, &i.SleepEnd,
			&i.DurationMinutes, &i.Quality, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const countSleepLogsByClient = `-- name: CountSleepLogsByClient :one
SELECT COUNT(*) FROM sleep_logs WHERE client_id = $1`

// CountSleepLogsByClient returns the total count of sleep logs for a client.
func (q *Queries) CountSleepLogsByClient(ctx context.Context, clientID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countSleepLogsByClient, clientID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
