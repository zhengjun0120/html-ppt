<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { usageApi, type UsageDay, type UsageOverview } from '@/api/usage'

const ov = ref<UsageOverview | null>(null)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    ov.value = await usageApi.overview()
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
})

// 柱高的基准是 30 天峰值：任何一天都和最高那天比，趋势形状一眼可读
const maxTotal = computed(() => Math.max(1, ...(ov.value?.daily.map((d) => d.total) ?? [0])))

const hasAny = computed(() => (ov.value?.daily ?? []).some((d) => d.total > 0))

function barH(d: UsageDay): string {
  const pct = Math.max(2, Math.round((d.total / maxTotal.value) * 100))
  return `${pct}%`
}

// 三段各占柱子的百分比；分母用 day.total（= 三段之和），柱内比例即真实构成
function segH(part: number, total: number): string {
  if (total <= 0) return '0%'
  return `${Math.round((part / total) * 100)}%`
}

function fmt(n: number): string {
  return n.toLocaleString()
}

function rateText(rate: number): string {
  return `${Math.round(rate * 100)}%`
}

// 悬浮卡片贴边翻转：贴着左右边缘的柱子，卡片朝卡片内侧展开，避免溢出容器
function tipPos(i: number): string {
  const n = ov.value?.daily.length ?? 0
  if (n === 0 || i <= 2) return 'left-0'
  if (i >= n - 3) return 'right-0'
  return 'left-1/2 -translate-x-1/2'
}

// 轴标签：首、每 7 天、末
function showAxis(i: number): boolean {
  return i === 0 || i === (ov.value?.daily.length ?? 0) - 1 || i % 7 === 0
}
</script>

<template>
  <div class="mx-auto flex max-w-[1100px] flex-col gap-4 p-6">
    <h1 class="text-[16px] font-bold">用量</h1>

    <div v-if="loading" class="rounded-card border border-line bg-surface p-8 text-center text-[13px] text-ink-3">
      加载中…
    </div>

    <div v-else-if="error" class="rounded-card border border-line bg-surface p-6 text-[13px] text-danger">
      {{ error }}
    </div>

    <template v-else-if="ov">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <section class="rounded-card border border-line bg-surface p-5">
          <h2 class="text-[13px] font-semibold text-ink-2">今日</h2>
          <p class="mt-2 font-mono text-[26px] font-bold leading-none">
            {{ fmt(ov.today.tokens) }}<span class="ml-1 text-[13px] font-normal text-ink-3">tok</span>
          </p>
          <p class="mt-2 text-[12.5px] text-ink-3">
            {{ ov.today.calls }} 次调用 · 缓存命中
            <span class="font-mono" title="缓存命中 token ÷ 输入 token；越高越省钱">{{ rateText(ov.today.cached_rate) }}</span>
          </p>
        </section>
        <section class="rounded-card border border-line bg-surface p-5">
          <h2 class="text-[13px] font-semibold text-ink-2">本月</h2>
          <p class="mt-2 font-mono text-[26px] font-bold leading-none">
            {{ fmt(ov.month.tokens) }}<span class="ml-1 text-[13px] font-normal text-ink-3">tok</span>
          </p>
          <p class="mt-2 text-[12.5px] text-ink-3">
            {{ ov.month.calls }} 次调用 · 缓存命中
            <span class="font-mono" title="缓存命中 token ÷ 输入 token；越高越省钱">{{ rateText(ov.month.cached_rate) }}</span>
          </p>
        </section>
      </div>

      <section class="rounded-card border border-line bg-surface p-5">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-[13px] font-semibold text-ink-2">近 30 天</h2>
          <div class="flex items-center gap-3 text-[11.5px] text-ink-3">
            <span class="inline-flex items-center gap-1"><i class="h-2 w-2 rounded-[2px] bg-success" />缓存命中</span>
            <span class="inline-flex items-center gap-1"><i class="h-2 w-2 rounded-[2px] bg-accent" />未命中输入</span>
            <span class="inline-flex items-center gap-1"><i class="h-2 w-2 rounded-[2px] bg-info" />输出</span>
          </div>
        </div>

        <div v-if="!hasAny" class="py-10 text-center text-[13px] text-ink-3">
          还没有用量记录——生成一份文稿或发起一次对话后，这里会出现你的用量。
        </div>

        <template v-else>
          <div class="flex h-44 items-end gap-[3px]">
            <div
              v-for="(d, i) in ov.daily"
              :key="d.date"
              class="group relative flex h-full min-w-0 flex-1 flex-col items-center justify-end"
            >
              <!-- 悬浮明细卡：整列都是热区（比细柱好悬停）；贴边列翻转对齐 -->
              <div
                v-if="d.total > 0"
                class="pointer-events-none absolute bottom-full z-10 mb-1 hidden w-max min-w-[150px] rounded-control border border-line bg-surface px-3 py-2 text-left shadow-md group-hover:block"
                :class="tipPos(i)"
              >
                <p class="font-mono text-[11px] text-ink-3">{{ d.date }}</p>
                <p class="mt-1 flex items-center gap-1.5 text-[11.5px]">
                  <i class="h-2 w-2 shrink-0 rounded-[2px] bg-success" />缓存命中
                  <span class="ml-auto pl-3 font-mono">{{ fmt(d.cached) }}</span>
                </p>
                <p class="mt-0.5 flex items-center gap-1.5 text-[11.5px]">
                  <i class="h-2 w-2 shrink-0 rounded-[2px] bg-accent" />未命中输入
                  <span class="ml-auto pl-3 font-mono">{{ fmt(d.uncached_in) }}</span>
                </p>
                <p class="mt-0.5 flex items-center gap-1.5 text-[11.5px]">
                  <i class="h-2 w-2 shrink-0 rounded-[2px] bg-info" />输出
                  <span class="ml-auto pl-3 font-mono">{{ fmt(d.output) }}</span>
                </p>
                <p class="mt-1 border-t border-line pt-1 text-[11px] text-ink-3">
                  共 {{ fmt(d.total) }} tok · {{ d.calls }} 次调用
                </p>
              </div>
              <div
                v-if="d.total > 0"
                class="flex w-full max-w-[20px] flex-col justify-end overflow-hidden rounded-t-[3px]"
                :style="{ height: barH(d) }"
              >
                <div class="w-full bg-info" :style="{ height: segH(d.output, d.total) }" />
                <div class="w-full bg-accent" :style="{ height: segH(d.uncached_in, d.total) }" />
                <div class="w-full bg-success" :style="{ height: segH(d.cached, d.total) }" />
              </div>
              <div v-else class="h-[2px] w-full max-w-[20px] rounded-full bg-line" :title="`${d.date.slice(5)}：无调用`" />
              <span
                class="mt-1 h-[14px] font-mono text-[10px] leading-none text-ink-3"
                :class="showAxis(i) ? '' : 'invisible'"
              >{{ d.date.slice(5) }}</span>
            </div>
          </div>
        </template>
      </section>
    </template>
  </div>
</template>
