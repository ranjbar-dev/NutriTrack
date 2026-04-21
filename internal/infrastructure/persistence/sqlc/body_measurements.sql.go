package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const upsertBodyMeasurement = `-- name: UpsertBodyMeasurement :one
INSERT INTO body_measurements (client_id, local_id, measured_at, measured_date, weight_kg, height_cm, waist_cm, hip_cm, chest_cm, arm_cm, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (client_id, local_id) DO UPDATE SET client_id = EXCLUDED.client_id
RETURNING id, client_id, local_id, measured_at, measured_date, weight_kg, height_cm, waist_cm, hip_cm, chest_cm, arm_cm, notes, created_at, (xmax::text::bigint = 0) AS inserted`

// UpsertBodyMeasurementParams holds parameters for upserting a body measurement.
type UpsertBodyMeasurementParams struct {
	ClientID     uuid.UUID       `db:"client_id"`
	LocalID      string          `db:"local_id"`
	MeasuredAt   time.Time       `db:"measured_at"`
	MeasuredDate time.Time       `db:"measured_date"`
	WeightKg     *pgtype.Numeric `db:"weight_kg"`
	HeightCm     *pgtype.Numeric `db:"height_cm"`
	WaistCm      *pgtype.Numeric `db:"waist_cm"`
	HipCm        *pgtype.Numeric `db:"hip_cm"`
	ChestCm      *pgtype.Numeric `db:"chest_cm"`
	ArmCm        *pgtype.Numeric `db:"arm_cm"`
	Notes        string          `db:"notes"`
}

// UpsertBodyMeasurement inserts or updates a body measurement and returns the row plus whether it was newly inserted.
func (q *Queries) UpsertBodyMeasurement(ctx context.Context, arg UpsertBodyMeasurementParams) (BodyMeasurement, bool, error) {
	row := q.db.QueryRow(ctx, upsertBodyMeasurement,
		arg.ClientID, arg.LocalID, arg.MeasuredAt, arg.MeasuredDate,
		arg.WeightKg, arg.HeightCm, arg.WaistCm, arg.HipCm, arg.ChestCm, arg.ArmCm, arg.Notes,
	)
	var i BodyMeasurement
	var inserted bool
	err := row.Scan(
		&i.ID, &i.ClientID, &i.LocalID, &i.MeasuredAt, &i.MeasuredDate,
		&i.WeightKg, &i.HeightCm, &i.WaistCm, &i.HipCm, &i.ChestCm, &i.ArmCm, &i.Notes, &i.CreatedAt,
		&inserted,
	)
	return i, inserted, err
}

// ListBodyMeasurementsByClientAndDateParams holds parameters for listing body measurements by client and date.
type ListBodyMeasurementsByClientAndDateParams struct {
	ClientID     uuid.UUID `db:"client_id"`
	MeasuredDate time.Time `db:"measured_date"`
}

const listBodyMeasurementsByClientAndDate = `-- name: ListBodyMeasurementsByClientAndDate :many
SELECT id, client_id, local_id, measured_at, measured_date, weight_kg, height_cm, waist_cm, hip_cm, chest_cm, arm_cm, notes, created_at
FROM body_measurements
WHERE client_id = $1 AND measured_date = $2
ORDER BY measured_at ASC`

// ListBodyMeasurementsByClientAndDate returns all body measurements for a client on a given date.
func (q *Queries) ListBodyMeasurementsByClientAndDate(ctx context.Context, arg ListBodyMeasurementsByClientAndDateParams) ([]BodyMeasurement, error) {
	rows, err := q.db.Query(ctx, listBodyMeasurementsByClientAndDate, arg.ClientID, arg.MeasuredDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BodyMeasurement
	for rows.Next() {
		var i BodyMeasurement
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.MeasuredAt, &i.MeasuredDate,
			&i.WeightKg, &i.HeightCm, &i.WaistCm, &i.HipCm, &i.ChestCm, &i.ArmCm, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ListBodyMeasurementsByClientParams holds parameters for paginated body measurement listing.
type ListBodyMeasurementsByClientParams struct {
	ClientID uuid.UUID `db:"client_id"`
	Limit    int32     `db:"limit"`
	Offset   int32     `db:"offset"`
}

const listBodyMeasurementsByClient = `-- name: ListBodyMeasurementsByClient :many
SELECT id, client_id, local_id, measured_at, measured_date, weight_kg, height_cm, waist_cm, hip_cm, chest_cm, arm_cm, notes, created_at
FROM body_measurements
WHERE client_id = $1
ORDER BY measured_date DESC, measured_at DESC
LIMIT $2 OFFSET $3`

// ListBodyMeasurementsByClient returns paginated body measurements for a client.
func (q *Queries) ListBodyMeasurementsByClient(ctx context.Context, arg ListBodyMeasurementsByClientParams) ([]BodyMeasurement, error) {
	rows, err := q.db.Query(ctx, listBodyMeasurementsByClient, arg.ClientID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BodyMeasurement
	for rows.Next() {
		var i BodyMeasurement
		if err := rows.Scan(
			&i.ID, &i.ClientID, &i.LocalID, &i.MeasuredAt, &i.MeasuredDate,
			&i.WeightKg, &i.HeightCm, &i.WaistCm, &i.HipCm, &i.ChestCm, &i.ArmCm, &i.Notes, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const countBodyMeasurementsByClient = `-- name: CountBodyMeasurementsByClient :one
SELECT COUNT(*) FROM body_measurements WHERE client_id = $1`

// CountBodyMeasurementsByClient returns the total count of body measurements for a client.
func (q *Queries) CountBodyMeasurementsByClient(ctx context.Context, clientID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countBodyMeasurementsByClient, clientID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
