<script setup lang="ts">
import { computed, ref } from 'vue'

import { PhPauseCircle, PhPaperPlaneRight } from '@phosphor-icons/vue'
import Button from '@/components/ui/Button.vue'
import Kbd from '@/components/ui/Kbd.vue'
import { useChatStore } from '@/stores/chat'

const props = defineProps<{ /** 向导锁（如选模板期间）：输入框禁用并说明原因 */ lockedHint?: string }>()
const chat = useChatStore()
const text = ref('')

const canSend = computed(() => chat.canSend && !props.lockedHint && text.value.trim() !== '')

function send() {
  if (!canSend.value) return
  const v = text.value
  text.value = ''
  void chat.send(v)
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
    <div class="flex items-end gap-2">
      <textarea
        v-model="text"
        rows="1"
        :disabled="chat.status === 'paused' || chat.restoring || !!lockedHint"
        :placeholder="lockedHint ?? (chat.restoring ? '正在载入上次对话…' : '想对这份文稿做什么？Enter 发送，Shift+Enter 换行')"
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
