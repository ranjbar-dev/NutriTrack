# 05-01 Summary: Database Migration & SQL Queries

## What Was Built

Created DB migration `000009_create_communication.up.sql` establishing:
- `food_request_status` enum (`pending`, `approved`, `rejected`)
- `messages` table: id, sender_id, receiver_id, content, attachment_path, attachment_name, attachment_mime, created_at, read_at
- `food_requests` table: id, client_id, nutritionist_id, food_name, description, status, reject_reason, reviewed_at, created_at
- 5 indexes: sender+receiver+created, receiver+read_at+created, food_requests client+created, food_requests nutritionist+status, food_requests client+status
- `LEAST/GREATEST` technique documented for conversation queries (sender_id < receiver_id canonical ordering)

SQL query files created:
- `backend/db/queries/messages.sql` — 6 queries: SendMessage, GetConversation, GetNewMessages, MarkRead, GetUnreadCount, GetMessageByID
- `backend/db/queries/food_requests.sql` — 6 queries: CreateFoodRequest, GetFoodRequestByID, ListClientFoodRequests, ListNutriPendingFoodRequests, ApproveFoodRequest, RejectFoodRequest

Also added to `backend/db/queries/users.sql`:
- `GetClientByIDForNutritionist` — verifies nutritionist owns client
- `UpdateClientProfile` — updates height, dob, gender (nutritionist-only fields)

## Decisions Made

- Used LEAST/GREATEST(sender_id, receiver_id) ordering for canonical conversation pair queries
- food_request_status defined as PostgreSQL enum (not text) for constraint enforcement at DB level
- Separate indexes for client-side and nutritionist-side food request queries
- Message attachment fields all nullable (text only is valid)

## Outcomes

- Migration file ready for `golang-migrate` execution
- SQL files ready for sqlc generation (hand-written since sqlc CLI not in CI)
