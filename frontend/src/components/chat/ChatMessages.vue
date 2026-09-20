<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'

import AskCard from './AskCard.vue'
import ToolCard from './ToolCard.vue'
import ToolGroupCard from './ToolGroupCard.vue'
import Empty from '@/components/ui/Empty.vue'
import { PhChatCircleDots } from '@phosphor-icons/vue'
import { renderMarkdown } from '@/lib/markdown'
import { useChatStore } from '@/stores/chat'
import type { ViewEvent, ViewTool } from '@/types/chat'

const chat = useChatStore()
const scrollBox = ref<HTMLElement | null>(null)

// 追加/流式增长时跟随滚动，但只在用户本来就贴着底部时：往上翻了就让它停住，
// 否则思考/长输出流式期间用户会被强制摁回底部，根本没法回看。
// （签名 = 事件数 + 最后一条文本长度）
const pinned = ref(true)
function onScroll() {
  const box = scrollBox.value
  if (box) pinned.value = box.scrollHeight - box.scrollTop - box.clientHeight < 80
}
const growthSig = () => {
  const last = chat.events[chat.events.length - 1]
  return `${chat.events.length}:${last && 'text' in last ? last.text.length : 0}`
}
watch(growthSig, () => {
  void nextTick(() => {
    const box = scrollBox.value
    if (box && pinned.value) box.scrollTop = box.scrollHeight
  })
})
// 切换/恢复会话：内容整体替换，重置为贴底
watch(
  () => chat.sessionId,
  () => {
    pinned.value = true
    void nextTick(() => {
      const box = scrollBox.value
      if (box) box.scrollTop = box.scrollHeight
    })
  },
)

// 思考行的预览：流式期间显示最新一行，让"思考中…"变成可读的进度
function thinkPreview(text: string): string {
  const lines = text.split('\n').map((l) => l.trim()).filter(Boolean)
  return lines.length ? lines[lines.length - 1] : ''
}

// —— 渲染项：把相邻同类型的工具调用收成一组（纯展示层，store 事件不动）——
// 一轮 run 里连读 8 个版式骨架这类"瀑布"收成一张"读取版式 ×8"的卡，
// 展开才看到每一笔。数组项 = 工具组，单个事件项 = 原样渲染。
type RenderItem = ViewEvent | ViewTool[]
const renderItems = computed<RenderItem[]>(() => {
  const out: RenderItem[] = []
  for (const e of chat.events) {
    if (e.kind !== 'tool') {
      out.push(e)
      continue
    }
    const last = out[out.length - 1]
    if (Array.isArray(last) && last[0].name === e.name) {
      last.push(e)
      continue
    }
    if (!Array.isArray(last) && last.kind === 'tool' && last.name === e.name) {
      out[out.length - 1] = [last, e]
      continue
    }
    out.push(e)
  }
  return out
})
function itemKey(item: RenderItem): string {
  return Array.isArray(item) ? `g-${item[0].id}` : String(item.id)
}
</script>

<template>
  <div ref="scrollBox" class="min-h-0 flex-1 overflow-y-auto p-3" @scroll.passive="onScroll">
    <Empty
      v-if="chat.events.length === 0"
      title="开始与 agent 对话"
      desc="描述你想要的演示文稿，agent 会逐页生成；也可以针对当前内容提修改意见。"
    >
      <template #icon><PhChatCircleDots /></template>
    </Empty>

    <div v-else class="flex flex-col gap-2.5">
      <template v-for="item in renderItems" :key="itemKey(item)">
        <!-- 相邻同类型工具的折叠组 -->
        <ToolGroupCard v-if="Array.isArray(item)" :name="item[0].name || ''" :tools="item" />

        <!-- 用户气泡：主色淡底（长文本压实色太吵） -->
        <div
          v-else-if="item.kind === 'user'"
          class="max-w-[88%] self-end rounded-xl rounded-br-sm border border-accent-border bg-accent-soft px-3 py-2 text-[13.5px] whitespace-pre-wrap break-words"
        >
          {{ item.text }}
        </div>

        <!-- agent 气泡 -->
        <div
          v-else-if="item.kind === 'agent'"
          class="max-w-[95%] self-start rounded-xl rounded-bl-sm bg-surface-2 px-3 py-2 text-[13.5px]"
        >
          <div class="prose-sm break-words [&_a]:text-accent [&_a]:underline [&_code]:rounded [&_code]:bg-code-bg [&_code]:px-1 [&_code]:font-mono [&_code]:text-[12px] [&_li]:ml-4 [&_li]:list-disc [&_ol]:list-decimal [&_p]:my-1 [&_p:first-child]:mt-0 [&_p:last-child]:mb-0 [&_ul]:list-disc [&_ul]:pl-1" v-html="renderMarkdown(item.text)" />
          <span v-if="item.streaming" class="mt-0.5 inline-block h-3.5 w-1.5 animate-pulse bg-accent align-middle" />
        </div>

        <!-- 思考块：折叠不占版面；流式期间折叠行显示最新一行思考（可读的进度感） -->
        <details v-else-if="item.kind === 'think'" class="max-w-[95%] self-start border-l-2 border-line-strong py-0.5 pl-2.5 text-[12px] text-ink-3">
          <summary class="block cursor-pointer select-none truncate hover:text-ink-2" :title="item.text">
            {{ item.streaming ? thinkPreview(item.text) || '思考中…' : '已展开思考过程' }}
          </summary>
          <div class="mt-1 max-h-72 overflow-y-auto whitespace-pre-wrap break-words">{{ item.text }}</div>
        </details>

        <!-- 子模型输出（视觉审查报告等）：默认展开边跑边读 -->
        <details
          v-else-if="item.kind === 'sub'"
          open
          class="max-w-[95%] self-start border-l-2 border-accent py-0.5 pl-2.5 text-[12px] text-ink-3"
        >
          <summary class="cursor-pointer select-none">子任务输出{{ item.streaming ? '（进行中）' : '' }}</summary>
          <div class="mt-1 max-h-72 overflow-y-auto whitespace-pre-wrap break-words font-mono text-[11.5px]">{{ item.text }}</div>
        </details>

        <!-- 单个工具卡（相邻同类已被收组） -->
        <ToolCard v-else-if="item.kind === 'tool'" :tool="item" />

        <!-- ask_user 提问卡 -->
        <AskCard
          v-else-if="item.kind === 'ask'"
          :tool-call-id="item.toolCallId"
          :questions="item.questions"
          :answered="item.answered"
          :note="item.note"
          :active="item.active"
        />

        <!-- 完成meta：弱化小字 -->
        <p v-else-if="item.kind === 'done'" class="self-center text-[10.5px] text-ink-3">
          本轮完成<template v-if="item.totalTokens"> · tokens {{ item.totalTokens.toLocaleString() }}<template v-if="item.cachedTokens">（缓存 {{ item.cachedTokens.toLocaleString() }}）</template></template>
        </p>

        <!-- 错误 -->
        <div v-else-if="item.kind === 'error'" class="max-w-[95%] self-start rounded-control bg-danger-soft px-3 py-2 text-[13px] text-danger">
          {{ item.text }}
        </div>

        <!-- 系统提示行 -->
        <p v-else-if="item.kind === 'sys'" class="self-center text-[10.5px] text-ink-3">{{ item.text }}</p>
      </template>
    </div>
  </div>
</template>
