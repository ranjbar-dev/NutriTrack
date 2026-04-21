package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const upsertExerciseLog = `-- name: UpsertExerciseLog :one
INSERT INTO exercise_logs (client_id, local_id, logged_at, logged_date, exercise_name, duration_minutes, calories_burned, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING id, client_id, local_id, logged_at, logged_date, exercise_name, duration_minutes, calories_burned, notes, created_at, (xmax::text::bigint = 0) AS inserted`

// UpsertExerciseLogParams holds parameters for upserting an exercise log.
type UpsertExerciseLogParams struct {
	ClientID        uuid.UUID `db:"client_id"`
	LocalID         string    `db:"local_id"`
	LoggedAt        time.Time `db:"logged_at"`
	LoggedDate      time.Time `db:"logged_date"`
	ExerciseName    string    `db:"exercise_name"`
	DurationMinutes int32     `db:"duration_minutes"`
	CaloriesBurned  int32     `db:"calories_burned"`
	Notes           string    `db:"notes"`
}

// UpsertExerciseLog inserts or updates an exercise log and returns the row plus whether it was newly inserted.
func (q *Queries) UpsertExerciseLog(ctx context.Context, arg UpsertExerciseLogParams) (ExerciseLog, bool, error) {
	row := q.db.QueryRow(ctx, upsertExerciseLog,
		arg.ClientID, arg.LocalID, arg.LoggedAt, arg.LoggedDate,
		arg.ExerciseName, arg.DurationMinutes, arg.CaloriesBurned, arg.Notes,
	)
	var i ExerciseLog
	var inserted bool
	err := row.Scan(
		&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate,
		&i.ExerciseName, &i.DurationMinutes, &i.CaloriesBurned, &i.Notes, &i.CreatedAt,
		&inserted,
	)
	return i, inserted, err
}

// ListExerciseLogsByClientAndDateParams holds parameters for listing exercise logs by client and date.
type ListExerciseLogsByClientAndDateParams struct {
	ClientID   uuid.UUID `db:"client_id"`
	LoggedDate time.Time `db:"logged_date"`
}

const listExerciseLogsByClientAndDate = `-- name: ListExerciseLogsByClientAndDate :many
SELECT id, client_id, local_id, logged_at, logged_date, exercise_name, duration_minutes, calories_burned, notes, created_at
FROM exercise_logs
WHERE client_id = $1 AND logged_date = $2
ORDER BY logged_at ASC`

// ListExerciseLogsByClientAndDate returns all exercise logs for a client on a given date.
func (q *Queries) ListExerciseLogsByClientAndDate(ctx context.Context, arg ListExerciseLogsByClientAndDateParams) ([]ExerciseLog, error) {
	rows, err := q.db.Query(ctx, listExerciseLogsByClientAndDate, arg.ClientID, arg.LoggedDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExerciseLog
	for rows.Next() {
		var i ExerciseLog
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate,
			&i.ExerciseName, &i.DurationMinutes, &i.CaloriesBurned, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ListExerciseLogsByClientParams holds parameters for paginated exercise log listing.
type ListExerciseLogsByClientParams struct {
	ClientID uuid.UUID `db:"client_id"`
	Limit    int32     `db:"limit"`
	Offset   int32     `db:"offset"`
}

const listExerciseLogsByClient = `-- name: ListExerciseLogsByClient :many
SELECT id, client_id, local_id, logged_at, logged_date, exercise_name, duration_minutes, calories_burned, notes, created_at
FROM exercise_logs
WHERE client_id = $1
ORDER BY logged_date DESC, logged_at DESC
LIMIT $2 OFFSET $3`

// ListExerciseLogsByClient returns paginated exercise logs for a client.
func (q *Queries) ListExerciseLogsByClient(ctx context.Context, arg ListExerciseLogsByClientParams) ([]ExerciseLog, error) {
	rows, err := q.db.Query(ctx, listExerciseLogsByClient, arg.ClientID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExerciseLog
	for rows.Next() {
		var i ExerciseLog
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate,
			&i.ExerciseName, &i.DurationMinutes, &i.CaloriesBurned, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const countExerciseLogsByClient = `-- name: CountExerciseLogsByClient :one
SELECT COUNT(*) FROM exercise_logs WHERE client_id = $1`

// CountExerciseLogsByClient returns the total count of exercise logs for a client.
func (q *Queries) CountExerciseLogsByClient(ctx context.Context, clientID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countExerciseLogsByClient, clientID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
