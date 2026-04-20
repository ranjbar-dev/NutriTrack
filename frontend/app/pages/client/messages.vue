<script setup lang="ts">
definePageMeta({
  layout: 'client',
  middleware: ['role-guard'],
  roles: ['client'],
})

const authStore = useAuthStore()
const messageStore = useMessageStore()

const myId = computed(() => authStore.user?.id ?? '')
const newContent = ref('')
const selectedFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const sending = ref(false)
const errorMsg = ref<string | null>(null)

const route = useRoute()
const partnerId = computed(() => route.query.partner as string | undefined)
const hasPartner = computed(() => Boolean(partnerId.value))

onMounted(async () => {
  await messageStore.fetchUnreadCount()
})

let pollInterval: ReturnType<typeof setInterval> | null = null
let lastSeen = new Date().toISOString()

watch(partnerId, async (pid) => {
  if (pollInterval) clearInterval(pollInterval)
  messageStore.clearMessages()
  if (!pid) return
  await messageStore.fetchMessages(pid)
  await messageStore.markRead(pid)
  lastSeen = messageStore.messages.at(-1)?.sent_at ?? new Date().toISOString()
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

    <div v-if="!hasPartner" class="rounded-2xl bg-white px-6 py-10 text-center shadow-sm">
      <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 text-2xl">
        💬
      </div>
      <p class="mt-4 font-bold text-gray-800">گفتگویی برای نمایش انتخاب نشده است</p>
      <p class="mt-2 text-sm text-gray-500">
        برای مشاهده پیام‌ها، از طریق پروفایل کارشناس تغذیه وارد گفتگو شوید.
      </p>
    </div>

    <template v-else>
      <div class="flex-1 overflow-y-auto rounded-2xl bg-white p-4 shadow-sm">
        <div v-if="messageStore.loading" class="space-y-3">
          <div v-for="index in 4" :key="index" class="h-16 animate-pulse rounded-2xl bg-gray-100" />
        </div>
        <div v-else-if="messageStore.messages.length === 0" class="py-16 text-center">
          <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 text-2xl">
            ✉️
          </div>
          <p class="mt-4 font-bold text-gray-800">هنوز پیامی رد و بدل نشده است</p>
          <p class="mt-2 text-sm text-gray-500">اولین پیام را برای کارشناس تغذیه خود ارسال کنید.</p>
        </div>
        <div v-else class="flex flex-col gap-3">
          <div
            v-for="msg in messageStore.messages"
            :key="msg.id"
            class="flex"
            :class="msg.sender_id === myId ? 'justify-start' : 'justify-end'"
          >
            <div
              class="max-w-[75%] rounded-2xl px-4 py-3 text-sm shadow-sm"
              :class="msg.sender_id === myId ? 'bg-emerald-100 text-emerald-900' : 'bg-gray-100 text-gray-800'"
            >
              <p v-if="msg.content">{{ msg.content }}</p>
              <div v-if="msg.attachment_type" class="mt-2 text-xs font-medium text-emerald-700">
                📎 {{ msg.attachment_name }}
              </div>
              <p class="mt-2 text-[10px] text-gray-400">{{ formatTime(msg.sent_at) }}</p>
            </div>
          </div>
        </div>
      </div>

      <div v-if="errorMsg" class="rounded-2xl border border-red-100 bg-red-50 px-4 py-3 text-sm text-red-600">
        {{ errorMsg }}
      </div>

      <div class="rounded-2xl bg-white p-3 shadow-sm">
        <p v-if="selectedFile" class="mb-3 rounded-xl bg-gray-50 px-3 py-2 text-xs text-gray-500">
          فایل انتخابی: {{ selectedFile.name }}
        </p>
        <div class="flex items-end gap-2">
        <textarea
          v-model="newContent"
          rows="2"
          placeholder="پیام خود را بنویسید..."
          class="min-h-[44px] flex-1 resize-none rounded-xl border border-gray-200 p-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-400"
          @keydown.enter.prevent="send"
        />
        <div class="flex flex-col gap-1">
          <button
            type="button"
            class="flex min-h-[44px] min-w-[44px] items-center justify-center rounded-xl text-gray-400 transition hover:bg-gray-50 hover:text-emerald-600"
            @click="fileInput?.click()"
          >
            📎
          </button>
          <input ref="fileInput" type="file" accept="image/jpeg,image/png,application/pdf" class="hidden" @change="onFileChange">
          <button
            type="button"
            :disabled="sending"
            class="min-h-[44px] rounded-xl bg-emerald-500 px-4 py-2 text-sm text-white disabled:opacity-50"
            @click="send"
          >
            ارسال
          </button>
        </div>
      </div>
      </div>
    </template>
  </div>
</template>
