<script setup lang="ts">
import { PhCaretDown, PhCaretUp, PhCheck, PhPlus, PhX } from '@phosphor-icons/vue'
import { computed, ref, watch } from 'vue'

import Button from '@/components/ui/Button.vue'
import { OutlineConflictError, type Outline, type OutlinePage } from '@/api/deckV2'
import { useWizardStore } from '@/stores/wizard'
import { useToast } from '@/stores/toast'

/**
 * 大纲面板（gate 1 的主区域）：结构化大纲的可编辑卡片。
 * 双通道之一——这里的编辑走 PUT 直改；另一通道（对话里让 agent 改）会通过
 * outline_updated 事件触发 wizard.refresh()，version 乐观锁保证两边不打架。
 */

const wizard = useWizardStore()
const toast = useToast()

/** 本地草稿（确认前可自由编辑；保存成功后 version 对齐服务端） */
const draft = ref<Outline | null>(null)
const dirty = computed(() => draft.value !== null)

const ROLE_LABELS: Record<string, string> = {
  cover: '封面', toc: '目录', divider: '章节', content: '正文',
  data: '数据', quote: '引用', code: '代码', cta: '行动', thanks: '收尾',
}

function startEdit() {
  if (!wizard.outline) return
  draft.value = JSON.parse(JSON.stringify(wizard.outline)) as Outline
}

function discard() {
  draft.value = null
}

function mutate(fn: (pages: OutlinePage[]) => void) {
  if (!draft.value) return
  fn(draft.value.pages)
  renumber()
}

function renumber() {
  if (!draft.value) return
  draft.value.pages.forEach((p, i) => (p.no = i + 1))
}

function removePage(i: number) {
  mutate((pages) => pages.splice(i, 1))
}

function move(i: number, delta: number) {
  mutate((pages) => {
    const j = i + delta
    if (j < 0 || j >= pages.length) return
    const t = pages[i]
    pages[i] = pages[j]
    pages[j] = t
  })
}

function addPage() {
  mutate((pages) => pages.push({ no: pages.length + 1, role: 'content', title: '新页面标题', points: [] }))
}

async function saveOnly(): Promise<boolean> {
  if (!draft.value) return true
  try {
    await wizard.saveOutline(draft.value)
    draft.value = null
    return true
  } catch (e) {
    if (e instanceof OutlineConflictError) {
      toast.error('大纲刚被另一通道更新，已加载最新版，请重新编辑')
    } else {
      toast.error(e instanceof Error ? e.message : '保存失败')
    }
    return false
  }
}

async function confirm() {
  if (!(await saveOnly())) return
  try {
    await wizard.confirmOutline()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '确认失败')
  }
}

// agent 通道更新了大纲（wizard.outline 变化）且本地不在编辑态 → 跟随
// （本地编辑态时以用户未保存的草稿为先，由 version 冲突提示兜底）
watch(
  () => wizard.outline?.version,
  () => {
    if (!dirty.value) draft.value = null
  },
)

// 首次渲染即进编辑态（面板本来就是编辑器）
if (wizard.outline && !draft.value) startEdit()
</script>

<template>
  <div class="mx-auto flex h-full w-full max-w-[860px] flex-col gap-3 overflow-y-auto p-5">
    <div class="flex items-baseline justify-between">
      <div>
        <h2 class="text-[16px] font-semibold text-ink">确认大纲</h2>
        <p class="mt-0.5 text-[12px] text-ink-3">
          直接编辑下面的卡片，或在右侧对话里让 agent 修改。确认后进入模板选择。
        </p>
      </div>
      <div class="flex items-center gap-2">
        <Button v-if="dirty" size="sm" :disabled="wizard.busy" @click="discard">
          <PhX :size="12" /> 放弃改动
        </Button>
        <Button v-if="dirty" size="sm" :disabled="wizard.busy" @click="saveOnly">
          <PhCheck :size="12" /> 保存
        </Button>
        <Button size="sm" variant="primary" :disabled="wizard.busy" @click="confirm">
          <PhCheck :size="12" /> 确认大纲
        </Button>
      </div>
    </div>

    <div v-if="wizard.outline" class="flex flex-col gap-2">
      <div class="rounded-control border border-line bg-surface px-3 py-2 text-[12.5px] text-ink-2">
        <span class="font-semibold text-ink">{{ wizard.outline.title }}</span>
        <span v-if="wizard.outline.meta?.audience" class="ml-2">· {{ wizard.outline.meta.audience }}</span>
        <span v-if="wizard.outline.meta?.duration_min" class="ml-2">· 约 {{ wizard.outline.meta.duration_min }} 分钟</span>
        <span class="ml-2 text-ink-3">v{{ dirty ? draft!.version : wizard.outline.version }}</span>
      </div>

      <template v-for="(p, i) in dirty ? draft!.pages : wizard.outline.pages" :key="i">
        <div class="rounded-control border border-line bg-surface p-3 transition-colors hover:border-line-strong">
          <div class="flex items-center gap-2">
            <span class="rounded bg-surface-2 px-1.5 py-0.5 font-mono text-[10.5px] text-ink-3">{{ p.no }}</span>
            <span class="rounded bg-accent-soft px-1.5 py-0.5 text-[10.5px] text-accent">{{ ROLE_LABELS[p.role] ?? p.role }}</span>
            <input
              v-model="p.title"
              class="min-w-0 flex-1 rounded border border-transparent bg-transparent px-1 py-0.5 text-[13px] font-medium text-ink outline-none hover:border-line focus-visible:border-accent"
              :readonly="!dirty"
              @focus="dirty || startEdit()"
            />
            <span v-if="dirty" class="flex shrink-0 items-center gap-0.5">
              <button class="cursor-pointer rounded p-1 text-ink-3 hover:bg-surface-2 hover:text-ink disabled:opacity-40" :disabled="i === 0" aria-label="上移" @click="move(i, -1)">
                <PhCaretUp :size="13" />
              </button>
              <button class="cursor-pointer rounded p-1 text-ink-3 hover:bg-surface-2 hover:text-ink disabled:opacity-40" :disabled="i === (draft?.pages.length ?? 0) - 1" aria-label="下移" @click="move(i, 1)">
                <PhCaretDown :size="13" />
              </button>
              <button class="cursor-pointer rounded p-1 text-ink-3 hover:bg-surface-2 hover:text-danger" aria-label="删除页" @click="removePage(i)">
                <PhX :size="13" />
              </button>
            </span>
          </div>
          <ul v-if="p.points?.length" class="mt-1.5 flex flex-col gap-0.5 pl-1 text-[12px] text-ink-2">
            <li v-for="(pt, j) in p.points" :key="j" class="flex gap-1.5">
              <span class="text-ink-3">·</span>
              <span
                contenteditable
                class="min-w-0 flex-1 rounded border border-transparent px-0.5 outline-none hover:border-line focus-visible:border-accent"
                @blur="dirty && (p.points![j] = ($event.target as HTMLElement).textContent?.trim() || pt)"
              >{{ pt }}</span>
            </li>
          </ul>
          <button
            v-if="dirty"
            class="mt-1.5 cursor-pointer text-[11px] text-ink-3 hover:text-accent"
            @click="mutate((pages) => { (pages[i].points ??= []).push('新要点') })"
          >
            <PhPlus :size="10" class="inline" /> 要点
          </button>
        </div>
      </template>

      <button
        v-if="dirty"
        class="flex cursor-pointer items-center justify-center gap-1 rounded-control border border-dashed border-line py-2 text-[12px] text-ink-3 transition-colors hover:border-accent hover:text-accent"
        @click="addPage"
      >
        <PhPlus :size="12" /> 加一页
      </button>
    </div>
  </div>
</template>
