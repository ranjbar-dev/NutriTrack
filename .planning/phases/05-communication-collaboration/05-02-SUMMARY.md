# 05-02 Summary: sqlc Go Files & Model Updates

## What Was Built

Hand-written sqlc Go files matching v1.30.0 output format:

**`backend/internal/repository/sqlc/messages.sql.go`** — 6 methods:
- `SendMessage`, `GetConversation`, `GetNewMessages`, `MarkRead`, `GetUnreadCount`, `GetMessageByID`

**`backend/internal/repository/sqlc/food_requests.sql.go`** — 6 methods:
- `CreateFoodRequest`, `GetFoodRequestByID`, `ListClientFoodRequests`, `ListNutriPendingFoodRequests`, `ApproveFoodRequest`, `RejectFoodRequest`

**`backend/internal/repository/sqlc/users.sql.go`** — appended 2 methods:
- `GetClientByIDForNutritionist`, `UpdateClientProfile`

**`backend/internal/repository/sqlc/models.go`** — added:
- `FoodRequestStatus` string enum with `Scan`/`Value` driver interface methods
- `NullFoodRequestStatus` struct
- `Message` struct (9 fields, nullable attachment fields as `pgtype.Text`)
- `FoodRequest` struct (9 fields, nullable description/reject_reason/reviewed_at)

**`backend/internal/repository/sqlc/querier.go`** — added 14 new method signatures

## Key Technical Details

- `pgtype.UUID.Bytes` is `[16]byte` — must convert with `uuid.UUID(x.Bytes).String()`
- `:many` queries use `q.db.Query` + `defer rows.Close()` + `items := []Type{}` (not nil)
- `:exec` queries return only error
- `:one` queries use `q.db.QueryRow`
- Nullable TEXT → `pgtype.Text`, Nullable UUID → `pgtype.UUID`, Nullable TIMESTAMPTZ → `pgtype.Timestamptz`

## Outcomes

- All sqlc Go files compile cleanly
- Querier interface complete with all 14 new signatures
