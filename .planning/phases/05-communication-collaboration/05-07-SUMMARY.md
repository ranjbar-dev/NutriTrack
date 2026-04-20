# 05-07 Summary: Frontend Pages

## What Was Built

### Client Pages

**`frontend/app/pages/client/messages.vue`**
- Conversation with assigned nutritionist
- Text input + optional file attachment (drag-and-drop / tap)
- Auto-starts polling via `useMessagePolling` composable
- Chronological message list with sent/received styling
- Read receipt indicator

**`frontend/app/pages/client/food-requests.vue`**
- Create new food request form (food_name + optional description)
- List of own requests with status badge (pending/approved/rejected)
- Shows reject_reason when rejected

### Nutritionist Pages

**`frontend/app/pages/nutritionist/messages/index.vue`**
- Lists all clients with unread badge counts
- Tap to open conversation

**`frontend/app/pages/nutritionist/messages/[partnerId].vue`**
- Full conversation view with client
- Same layout as client side; nutritionist can also send attachments
- Polling active while page is open

**`frontend/app/pages/nutritionist/food-requests.vue`**
- Lists pending food requests from clients
- Approve button → redirects to `/nutritionist/foods/create?name=<food_name>`
- Reject button → inline reason input + confirm

**`frontend/app/pages/nutritionist/clients.vue`** (replaced stub)
- Search bar (name/mobile)
- Active/inactive filter tabs
- Sort by name or last activity
- Client cards linking to profile

**`frontend/app/pages/nutritionist/clients/[clientId]/index.vue`**
- Full client profile: personal info (height, DOB, gender — editable by nutritionist)
- Quick action buttons: New Plan, Send Message, Activate/Deactivate
- Tabs: Overview, Tracking, Plans
- Edit mode toggled inline

## Key Technical Details

- All pages use `<script setup>` with Composition API
- `useMessagePolling` composable starts/stops automatically via lifecycle hooks
- File inputs in messages pages accept `image/jpeg,image/png,application/pdf`
- Client profile edit updates via `clientManagementStore.updateClientProfile`

## Outcomes

- All 7 frontend pages created
- Phase 5 frontend implementation complete
