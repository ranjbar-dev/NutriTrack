<script setup lang="ts">
definePageMeta({ layout: 'client' })

const authStore = useAuthStore()
const messageStore = useMessageStore()
const { useShamsiDate } = useShamsiDate ? { useShamsiDate } : { useShamsiDate: null }

// The client's partner is their nutritionist
// We load messages and poll for new ones
const myId = computed(() => authStore.user?.id ?? '')

// We need to find the nutritionist ID — fetch it from the profile or the active plan
// For simplicity, messages show "پیام‌های من" and we send using receiver_id
const newContent = ref('')
const selectedFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const sending = ref(false)
const errorMsg = ref<string | null>(null)

// Determine partner — stored in the auth user profile from GET /auth/me
// We'll use a query param ?partner=<id> if available, otherwise show all messages
const route = useRoute()
const partnerId = computed(() => route.query.partner as string | undefined)

onMounted(async () => {
  await messageStore.fetchUnreadCount()
  if (partnerId.value) {
    await messageStore.fetchMessages(partnerId.value)
    await messageStore.markRead(partnerId.value)
  }
})

let pollInterval: ReturnType<typeof setInterval> | null = null
let lastSeen = new Date().toISOString()

watch(partnerId, (pid) => {
  if (pollInterval) clearInterval(pollInterval)
  if (!pid) return
  pollInterval = setInterval(async () => {
    const newMsgs = await messageStore.pollNewMessages(pid, lastSeen)
    if (newMsgs.length > 0) lastSeen = newMsgs[newMsgs.length - 1].sent_at
  }, 10_000)
}, { immediate: true })

onUnmounted(() => { if (pollInterval) clearInterval(pollInterval) })

async function send() {
  if (!partnerId.value) return
  if (!newContent.value.trim() && !selectedFile.value) return

  sending.value = true
  errorMsg.value = null
  await messageStore.sendMessage(
    partnerId.value,
    newContent.value.trim() || undefined,
    selectedFile.value ?? undefined,
  )
  if (!messageStore.error) {
    newContent.value = ''
    selectedFile.value = null
  }
  else {
    errorMsg.value = messageStore.error
  }
  sending.value = false
}

function onFileChange(e: Event) {
  const f = (e.target as HTMLInputElement).files?.[0]
  selectedFile.value = f ?? null
}

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('fa-IR', { hour: '2-digit', minute: '2-digit' })
}
</script>

<template>
  <div class="flex flex-col h-full p-4 gap-4">
    <h1 class="text-xl font-bold text-gray-800">پیام‌ها</h1>

    <div v-if="!partnerId" class="bg-white rounded-xl p-6 shadow-sm text-center text-gray-500">
      برای مشاهده پیام‌ها، از طریق پروفایل کارشناس تغذیه وارد گفتگو شوید.
    </div>

    <template v-else>
      <div class="flex-1 bg-white rounded-xl shadow-sm overflow-y-auto p-4 flex flex-col gap-3">
        <div v-if="messageStore.loading" class="text-center text-gray-400 text-sm">در حال بارگذاری...</div>
        <div v-else-if="messageStore.messages.length === 0" class="text-center text-gray-400 text-sm">
          هنوز پیامی ارسال نشده است.
        </div>
        <div
          v-for="msg in messageStore.messages"
          :key="msg.id"
          class="flex"
          :class="msg.sender_id === myId ? 'justify-start' : 'justify-end'"
        >
          <div
            class="max-w-[75%] rounded-2xl px-4 py-2 text-sm"
            :class="msg.sender_id === myId ? 'bg-emerald-100 text-emerald-900' : 'bg-gray-100 text-gray-800'"
          >
            <p v-if="msg.content">{{ msg.content }}</p>
            <div v-if="msg.attachment_type" class="mt-1 text-xs text-blue-600">
              📎 {{ msg.attachment_name }}
            </div>
            <p class="text-[10px] text-gray-400 mt-1">{{ formatTime(msg.sent_at) }}</p>
          </div>
        </div>
      </div>

      <div v-if="errorMsg" class="text-red-500 text-sm">{{ errorMsg }}</div>

      <div class="bg-white rounded-xl shadow-sm p-3 flex gap-2 items-end">
        <textarea
          v-model="newContent"
          rows="2"
          placeholder="پیام خود را بنویسید..."
          class="flex-1 resize-none rounded-lg border border-gray-200 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-400"
          @keydown.enter.prevent="send"
        />
        <div class="flex flex-col gap-1">
          <button
            class="text-gray-400 hover:text-emerald-600 transition"
            @click="fileInput?.click()"
          >
            📎
          </button>
          <input ref="fileInput" type="file" accept="image/jpeg,image/png,application/pdf" class="hidden" @change="onFileChange" />
          <button
            :disabled="sending"
            class="bg-emerald-500 text-white rounded-lg px-3 py-1 text-sm disabled:opacity-50"
            @click="send"
          >
            ارسال
          </button>
        </div>
      </div>
      <p v-if="selectedFile" class="text-xs text-gray-500">فایل انتخابی: {{ selectedFile.name }}</p>
    </template>
  </div>
</template>
