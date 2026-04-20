# 05-06 Summary: Frontend Types, Stores & Composable

## What Was Built

### Type Files

**`frontend/app/types/message.types.ts`**
- `Message` — id, sender_id, receiver_id, content, attachment fields, created_at, read_at
- `SendMessagePayload` — FormData-based for multipart upload
- `PollNewMessagesParams` — partnerId, since timestamp

**`frontend/app/types/foodRequest.types.ts`**
- `FoodRequest` — id, client_id, nutritionist_id, food_name, description, status (pending/approved/rejected), reject_reason, reviewed_at, created_at
- `CreateFoodRequestPayload`, `RejectFoodRequestPayload`

**`frontend/app/types/clientManagement.types.ts`**
- `ClientListItem` — condensed for list views (name, mobile, status, plan status, last activity)
- `ClientProfile` — full profile with all fields
- `UpdateClientProfilePayload` — height_cm, date_of_birth, gender

### Pinia Stores

**`frontend/app/stores/message.ts`**
- `conversations` map (partnerId → Message[])
- `unreadCount` ref
- Actions: sendMessage (FormData), loadConversation, pollNew, markRead, fetchUnreadCount

**`frontend/app/stores/foodRequest.ts`**
- `clientRequests` and `nutriRequests` arrays
- Actions: createRequest (client), listClientRequests, listNutriRequests, approveRequest, rejectRequest

**`frontend/app/stores/clientManagement.ts`**
- `clients` array, `selectedClient` ref, `searchQuery`/`filterStatus`/`sortBy` refs
- Actions: fetchClients, fetchClientProfile, activateClient, deactivateClient, updateClientProfile

### Composable

**`frontend/app/composables/useMessagePolling.ts`**
- Starts 10-second interval polling via `messageStore.pollNew(partnerId, lastTimestamp)`
- `onUnmounted` clears interval
- Exposes `startPolling(partnerId)`, `stopPolling()`

## Key Technical Details

- All stores use Pinia setup function style (`defineStore('name', () => {...})`)
- `useApi()` composable with `apiFetch<T>()` — credentials:'include', auto-refresh on 401
- Nuxt 4 auto-imports — no explicit imports needed for composables/stores in components
- `sendMessage` uses `FormData` for multipart compatibility

## Outcomes

- All type contracts, stores, and composable created
