<script setup lang="ts">
definePageMeta({ layout: 'nutritionist' })

const authStore = useAuthStore()
const messageStore = useMessageStore()
const route = useRoute()
const partnerId = route.params.partnerId as string

const myId = computed(() => authStore.user?.id ?? '')
const newContent = ref('')
const selectedFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const sending = ref(false)
const errorMsg = ref<string | null>(null)

let lastSeen = new Date().toISOString()

onMounted(async () => {
  await messageStore.fetchMessages(partnerId)
  await messageStore.markRead(partnerId)
  if (messageStore.messages.length > 0) {
    lastSeen = messageStore.messages[messageStore.messages.length - 1].sent_at
  }
})

const pollInterval = setInterval(async () => {
  const newMsgs = await messageStore.pollNewMessages(partnerId, lastSeen)
  if (newMsgs.length > 0) {
    lastSeen = newMsgs[newMsgs.length - 1].sent_at
    await messageStore.markRead(partnerId)
  }
}, 10_000)

onUnmounted(() => clearInterval(pollInterval))

async function send() {
  if (!newContent.value.trim() && !selectedFile.value) return
  sending.value = true
  errorMsg.value = null
  const msg = await messageStore.sendMessage(
    partnerId,
    newContent.value.trim() || undefined,
    selectedFile.value ?? undefined,
  )
  if (msg) {
    newContent.value = ''
    selectedFile.value = null
  } else {
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
    <div class="flex items-center gap-2">
      <NuxtLink to="/nutritionist/messages" class="text-gray-400 hover:text-gray-700">←</NuxtLink>
      <h1 class="text-xl font-bold text-gray-800">گفتگو</h1>
    </div>

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
          :class="msg.sender_id === myId ? 'bg-blue-100 text-blue-900' : 'bg-gray-100 text-gray-800'"
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
        class="flex-1 resize-none rounded-lg border border-gray-200 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
        @keydown.enter.prevent="send"
      />
      <div class="flex flex-col gap-1">
        <button class="text-gray-400 hover:text-blue-600 transition" @click="fileInput?.click()">📎</button>
        <input ref="fileInput" type="file" accept="image/jpeg,image/png,application/pdf" class="hidden" @change="onFileChange" />
        <button
          :disabled="sending"
          class="bg-blue-500 text-white rounded-lg px-3 py-1 text-sm disabled:opacity-50"
          @click="send"
        >
          ارسال
        </button>
      </div>
    </div>
    <p v-if="selectedFile" class="text-xs text-gray-500">فایل انتخابی: {{ selectedFile.name }}</p>
  </div>
</template>
