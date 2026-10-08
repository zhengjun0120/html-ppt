<script setup lang="ts">
import { computed, ref } from 'vue'

import { PhImage, PhPauseCircle, PhPaperPlaneRight, PhX, PhCircleNotch } from '@phosphor-icons/vue'
import Button from '@/components/ui/Button.vue'
import Kbd from '@/components/ui/Kbd.vue'
import { MAX_CHAT_IMAGES, fileToChatImage } from '@/lib/image'
import { useChatStore } from '@/stores/chat'
import { useToast } from '@/stores/toast'

const props = defineProps<{ /** 向导锁（如选模板期间）：输入框禁用并说明原因 */ lockedHint?: string }>()
const chat = useChatStore()
const toast = useToast()
const text = ref('')

// —— 用户附图（参考图/截图；压缩后随消息发，服务端 chatimg 校验 ≤3 张）——//
const pendingImages = ref<string[]>([])
const attaching = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

async function onPickImages(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (!files?.length) return
  attaching.value = true
  try {
    for (const f of Array.from(files)) {
      if (pendingImages.value.length >= MAX_CHAT_IMAGES) {
        toast.error(`一次最多附 ${MAX_CHAT_IMAGES} 张图`)
        break
      }
      try {
        pendingImages.value.push(await fileToChatImage(f))
      } catch (err) {
        toast.error(err instanceof Error ? err.message : '图片处理失败')
      }
    }
  } finally {
    attaching.value = false
    e.target.value = '' // 允许重选同一个文件
  }
}

function removePendingImage(i: number) {
  pendingImages.value.splice(i, 1)
}

const canSend = computed(
  () => chat.canSend && !props.lockedHint && !attaching.value && (text.value.trim() !== '' || pendingImages.value.length > 0),
)

function send() {
  if (!canSend.value) return
  const v = text.value
  const imgs = [...pendingImages.value]
  text.value = ''
  pendingImages.value = []
  void chat.send(v, false, imgs)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
}

function autoGrow(e: Event) {
  const el = e.target as HTMLTextAreaElement
  el.style.height = 'auto'
  el.style.height = `${Math.min(el.scrollHeight, 120)}px`
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <!-- 暂停态提示：硬拦截是后端闸门，前端必须把"该做什么"讲清楚 -->
    <p v-if="chat.status === 'paused'" class="text-[11.5px] text-warning">
      上一条提问还没有回答：请先在上方作答，或点「跳过」
    </p>
    <p v-else-if="lockedHint" class="text-[11.5px] text-warning">{{ lockedHint }}</p>
    <!-- 待发附图条 -->
    <div v-if="pendingImages.length" class="flex flex-wrap gap-1.5">
      <div v-for="(src, i) in pendingImages" :key="src" class="relative">
        <img :src="src" class="h-14 rounded-control border border-line" alt="待发送附图" />
        <button
          class="absolute -right-1.5 -top-1.5 flex h-4.5 w-4.5 items-center justify-center rounded-full bg-ink text-surface hover:bg-danger"
          :aria-label="`移除第 ${i + 1} 张附图`"
          @click="removePendingImage(i)"
        >
          <PhX :size="10" weight="bold" />
        </button>
      </div>
    </div>
    <div class="flex items-end gap-2">
      <input
        ref="fileInput"
        type="file"
        accept="image/png,image/jpeg,image/webp"
        multiple
        class="hidden"
        @change="onPickImages"
      />
      <button
        class="flex h-9 w-9 shrink-0 items-center justify-center rounded-control border border-line text-ink-2 hover:border-accent hover:text-accent disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="chat.status !== 'idle' || attaching || pendingImages.length >= MAX_CHAT_IMAGES"
        :title="pendingImages.length >= MAX_CHAT_IMAGES ? `一次最多附 ${MAX_CHAT_IMAGES} 张图` : '附图（最多 3 张，自动压缩）'"
        aria-label="附图"
        @click="fileInput?.click()"
      >
        <PhCircleNotch v-if="attaching" :size="15" class="animate-spin" />
        <PhImage v-else :size="15" />
      </button>
      <textarea
        v-model="text"
        rows="1"
        :disabled="chat.status === 'paused' || chat.restoring || !!lockedHint"
        :placeholder="lockedHint ?? (chat.restoring ? '正在载入上次对话…' : '想对这份文稿做什么？可附参考图。Enter 发送，Shift+Enter 换行')"
        class="max-h-[120px] w-full resize-none rounded-control border border-line bg-surface-2 px-3 py-2 text-[13.5px] text-ink outline-none transition-[border-color,box-shadow] placeholder:text-ink-3 focus-visible:border-accent focus-visible:shadow-[0_0_0_3px_var(--ring)] disabled:cursor-not-allowed disabled:opacity-55"
        @keydown="onKeydown"
        @input="autoGrow"
      />
      <Button
        v-if="chat.status !== 'streaming'"
        variant="primary"
        :disabled="!canSend"
        class="h-9 w-10 !px-0"
        aria-label="发送"
        @click="send"
      >
        <PhPaperPlaneRight :size="16" />
      </Button>
      <Button v-else variant="danger" class="h-9 w-10 !px-0" aria-label="停止生成" @click="chat.abort()">
        <PhPauseCircle :size="17" />
      </Button>
    </div>
    <p class="hidden items-center gap-1 text-[10.5px] text-ink-3 md:flex">
      <Kbd>Enter</Kbd> 发送 · <Kbd>Shift</Kbd>+<Kbd>Enter</Kbd> 换行
    </p>
  </div>
</template>
