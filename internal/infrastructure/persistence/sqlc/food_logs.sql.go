package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const upsertFoodLog = `-- name: UpsertFoodLog :one
INSERT INTO food_logs (client_id, local_id, logged_at, logged_date, food_id, food_name, quantity, unit, calories, protein, carbs, fat, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING id, client_id, local_id, logged_at, logged_date, food_id, food_name, quantity, unit, calories, protein, carbs, fat, notes, created_at, (xmax::text::bigint = 0) AS inserted`

// UpsertFoodLogParams holds parameters for upserting a food log.
type UpsertFoodLogParams struct {
	ClientID   uuid.UUID      `db:"client_id"`
	LocalID    string         `db:"local_id"`
	LoggedAt   time.Time      `db:"logged_at"`
	LoggedDate time.Time      `db:"logged_date"`
	FoodID     *uuid.UUID     `db:"food_id"`
	FoodName   string         `db:"food_name"`
	Quantity   pgtype.Numeric `db:"quantity"`
	Unit       string         `db:"unit"`
	Calories   pgtype.Numeric `db:"calories"`
	Protein    pgtype.Numeric `db:"protein"`
	Carbs      pgtype.Numeric `db:"carbs"`
	Fat        pgtype.Numeric `db:"fat"`
	Notes      string         `db:"notes"`
}

// UpsertFoodLog inserts or updates a food log and returns the row plus whether it was newly inserted.
func (q *Queries) UpsertFoodLog(ctx context.Context, arg UpsertFoodLogParams) (FoodLog, bool, error) {
	row := q.db.QueryRow(ctx, upsertFoodLog,
		arg.ClientID, arg.LocalID, arg.LoggedAt, arg.LoggedDate, arg.FoodID,
		arg.FoodName, arg.Quantity, arg.Unit, arg.Calories, arg.Protein, arg.Carbs, arg.Fat, arg.Notes,
	)
	var i FoodLog
	var inserted bool
	err := row.Scan(
		&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate, &i.FoodID,
		&i.FoodName, &i.Quantity, &i.Unit, &i.Calories, &i.Protein, &i.Carbs, &i.Fat, &i.Notes, &i.CreatedAt,
		&inserted,
	)
	return i, inserted, err
}

// ListFoodLogsByClientAndDateParams holds parameters for listing food logs by client and date.
type ListFoodLogsByClientAndDateParams struct {
	ClientID   uuid.UUID `db:"client_id"`
	LoggedDate time.Time `db:"logged_date"`
}

const listFoodLogsByClientAndDate = `-- name: ListFoodLogsByClientAndDate :many
SELECT id, client_id, local_id, logged_at, logged_date, food_id, food_name, quantity, unit, calories, protein, carbs, fat, notes, created_at
FROM food_logs
WHERE client_id = $1 AND logged_date = $2
ORDER BY logged_at ASC`

// ListFoodLogsByClientAndDate returns all food logs for a client on a given date.
func (q *Queries) ListFoodLogsByClientAndDate(ctx context.Context, arg ListFoodLogsByClientAndDateParams) ([]FoodLog, error) {
	rows, err := q.db.Query(ctx, listFoodLogsByClientAndDate, arg.ClientID, arg.LoggedDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FoodLog
	for rows.Next() {
		var i FoodLog
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate, &i.FoodID,
			&i.FoodName, &i.Quantity, &i.Unit, &i.Calories, &i.Protein, &i.Carbs, &i.Fat, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ListFoodLogsByClientParams holds parameters for paginated food log listing.
type ListFoodLogsByClientParams struct {
	ClientID uuid.UUID `db:"client_id"`
	Limit    int32     `db:"limit"`
	Offset   int32     `db:"offset"`
}

const listFoodLogsByClient = `-- name: ListFoodLogsByClient :many
SELECT id, client_id, local_id, logged_at, logged_date, food_id, food_name, quantity, unit, calories, protein, carbs, fat, notes, created_at
FROM food_logs
WHERE client_id = $1
ORDER BY logged_date DESC, logged_at DESC
LIMIT $2 OFFSET $3`

// ListFoodLogsByClient returns paginated food logs for a client.
func (q *Queries) ListFoodLogsByClient(ctx context.Context, arg ListFoodLogsByClientParams) ([]FoodLog, error) {
	rows, err := q.db.Query(ctx, listFoodLogsByClient, arg.ClientID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FoodLog
	for rows.Next() {
		var i FoodLog
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.LoggedAt, &i.LoggedDate, &i.FoodID,
			&i.FoodName, &i.Quantity, &i.Unit, &i.Calories, &i.Protein, &i.Carbs, &i.Fat, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const countFoodLogsByClient = `-- name: CountFoodLogsByClient :one
SELECT COUNT(*) FROM food_logs WHERE client_id = $1`

// CountFoodLogsByClient returns the total count of food logs for a client.
func (q *Queries) CountFoodLogsByClient(ctx context.Context, clientID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countFoodLogsByClient, clientID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
