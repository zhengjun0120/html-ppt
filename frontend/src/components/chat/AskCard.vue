<script setup lang="ts">
import { reactive, ref } from 'vue'

import type { AskQuestion } from '@/types/chat'
import Button from '@/components/ui/Button.vue'
import { useChatStore } from '@/stores/chat'

const props = defineProps<{
  toolCallId: string
  questions: AskQuestion[]
  answered: { question: string; answer: string }[] | null
  note: string | null
  active: boolean
}>()

const chat = useChatStore()
const submitting = ref(false)
// 每个问题：选中的选项 or 自定义输入
const picked = reactive<Record<number, string>>({})
const custom = reactive<Record<number, string>>({})

async function submit(skip: boolean) {
  if (submitting.value) return
  const answers: { question: string; answer: string }[] = []
  if (!skip) {
    for (const [i, q] of props.questions.entries()) {
      const answer = custom[i]?.trim() || picked[i]
      if (answer) answers.push({ question: q.question, answer })
    }
  }
  submitting.value = true
  try {
    await chat.answer(answers, skip ? '用户未作答，请使用推荐答案' : '用户未作答，请使用推荐答案')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="w-[95%] self-start rounded-[10px] border border-accent-border border-l-[3px] border-l-accent bg-surface px-3.5 py-3 text-[13px]">
    <div v-for="(q, i) in questions" :key="i" class="mb-3 last:mb-0">
      <p class="mb-2 font-semibold">{{ i + 1 }}. {{ q.question }}</p>

      <!-- 已回答：只读展示（回放/提交后） -->
      <template v-if="!active">
        <p class="rounded bg-surface-2 px-2.5 py-1.5 text-ink-2">
          {{ answered?.find((a) => a.question === q.question)?.answer ?? note ?? '（未作答）' }}
        </p>
      </template>

      <!-- 待答：选项 + 自定义输入 -->
      <template v-else>
        <div v-if="q.options?.length" class="mb-2 flex flex-wrap gap-2">
          <button
            v-for="opt in q.options"
            :key="opt"
            class="cursor-pointer rounded-full border px-3 py-1 text-[12.5px] transition-colors"
            :class="picked[i] === opt ? 'border-accent bg-accent-soft text-ink' : 'border-line bg-surface-2 text-ink-2 hover:border-accent'"
            @click="picked[i] = picked[i] === opt ? '' : opt"
          >
            {{ opt }}
          </button>
        </div>
        <input
          v-model="custom[i]"
          placeholder="或输入自定义回答…"
          class="w-full rounded-control border border-line bg-surface-2 px-2.5 py-1.5 text-[12.5px] text-ink outline-none placeholder:text-ink-3 focus-visible:border-accent focus-visible:shadow-[0_0_0_3px_var(--ring)]"
        />
      </template>
    </div>

    <div v-if="active" class="flex gap-2">
      <Button variant="primary" size="sm" :loading="submitting" @click="submit(false)">提交回答</Button>
      <Button variant="ghost" size="sm" :disabled="submitting" @click="submit(true)">跳过（用推荐答案）</Button>
    </div>
  </div>
</template>
