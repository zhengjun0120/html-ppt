<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'

import AskCard from './AskCard.vue'
import ToolCard from './ToolCard.vue'
import Empty from '@/components/ui/Empty.vue'
import { PhChatCircleDots } from '@phosphor-icons/vue'
import { renderMarkdown } from '@/lib/markdown'
import { useChatStore } from '@/stores/chat'

const chat = useChatStore()
const scrollBox = ref<HTMLElement | null>(null)

// 追加/流式增长时贴底（签名 = 事件数 + 最后一条文本长度）
const growthSig = () => {
  const last = chat.events[chat.events.length - 1]
  return `${chat.events.length}:${last && 'text' in last ? last.text.length : 0}`
}
watch(growthSig, () => {
  void nextTick(() => {
    const box = scrollBox.value
    if (box) box.scrollTop = box.scrollHeight
  })
})
</script>

<template>
  <div ref="scrollBox" class="min-h-0 flex-1 overflow-y-auto p-3">
    <Empty
      v-if="chat.events.length === 0"
      title="开始与 agent 对话"
      desc="描述你想要的演示文稿，agent 会逐页生成；也可以针对当前内容提修改意见。"
    >
      <template #icon><PhChatCircleDots /></template>
    </Empty>

    <div v-else class="flex flex-col gap-2.5">
      <template v-for="e in chat.events" :key="e.id">
        <!-- 用户气泡：主色淡底（长文本压实色太吵） -->
        <div
          v-if="e.kind === 'user'"
          class="max-w-[88%] self-end rounded-xl rounded-br-sm border border-accent-border bg-accent-soft px-3 py-2 text-[13.5px] whitespace-pre-wrap break-words"
        >
          {{ e.text }}
        </div>

        <!-- agent 气泡 -->
        <div
          v-else-if="e.kind === 'agent'"
          class="max-w-[95%] self-start rounded-xl rounded-bl-sm bg-surface-2 px-3 py-2 text-[13.5px]"
        >
          <div class="prose-sm break-words [&_a]:text-accent [&_a]:underline [&_code]:rounded [&_code]:bg-code-bg [&_code]:px-1 [&_code]:font-mono [&_code]:text-[12px] [&_li]:ml-4 [&_li]:list-disc [&_ol]:list-decimal [&_p]:my-1 [&_p:first-child]:mt-0 [&_p:last-child]:mb-0 [&_ul]:list-disc [&_ul]:pl-1" v-html="renderMarkdown(e.text)" />
          <span v-if="e.streaming" class="mt-0.5 inline-block h-3.5 w-1.5 animate-pulse bg-accent align-middle" />
        </div>

        <!-- 思考块：折叠，闭合时不参与布局 -->
        <details v-else-if="e.kind === 'think'" class="max-w-[95%] self-start border-l-2 border-line-strong py-0.5 pl-2.5 text-[12px] text-ink-3">
          <summary class="cursor-pointer select-none hover:text-ink-2">
            {{ e.streaming ? '思考中…' : '已展开思考过程' }}
          </summary>
          <div class="mt-1 max-h-72 overflow-y-auto whitespace-pre-wrap break-words">{{ e.text }}</div>
        </details>

        <!-- 子模型输出（视觉审查报告等）：默认展开边跑边读 -->
        <details
          v-else-if="e.kind === 'sub'"
          open
          class="max-w-[95%] self-start border-l-2 border-accent py-0.5 pl-2.5 text-[12px] text-ink-3"
        >
          <summary class="cursor-pointer select-none">子任务输出{{ e.streaming ? '（进行中）' : '' }}</summary>
          <div class="mt-1 max-h-72 overflow-y-auto whitespace-pre-wrap break-words font-mono text-[11.5px]">{{ e.text }}</div>
        </details>

        <!-- 工具卡 -->
        <ToolCard v-else-if="e.kind === 'tool'" :tool="e" />

        <!-- ask_user 提问卡 -->
        <AskCard
          v-else-if="e.kind === 'ask'"
          :tool-call-id="e.toolCallId"
          :questions="e.questions"
          :answered="e.answered"
          :note="e.note"
          :active="e.active"
        />

        <!-- 完成meta：弱化小字 -->
        <p v-else-if="e.kind === 'done'" class="self-center text-[10.5px] text-ink-3">
          本轮完成<template v-if="e.totalTokens"> · tokens {{ e.totalTokens.toLocaleString() }}<template v-if="e.cachedTokens">（缓存 {{ e.cachedTokens.toLocaleString() }}）</template></template>
        </p>

        <!-- 错误 -->
        <div v-else-if="e.kind === 'error'" class="max-w-[95%] self-start rounded-control bg-danger-soft px-3 py-2 text-[13px] text-danger">
          {{ e.text }}
        </div>

        <!-- 系统提示行 -->
        <p v-else-if="e.kind === 'sys'" class="self-center text-[10.5px] text-ink-3">{{ e.text }}</p>
      </template>
    </div>
  </div>
</template>
