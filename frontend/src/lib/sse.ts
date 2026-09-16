/**
 * POST + SSE 解析器。
 * 后端的 SSE 帧是 `event: <Type>\ndata: <JSON>\n\n`（LF 分隔、空行分帧）；
 * EventSource 不支持 POST 和 Authorization 头，所以用 fetch + ReadableStream 手工解析。
 */
export interface SseFrame {
  event: string
  data: string
}

function parseFrame(raw: string): SseFrame {
  let event = 'message'
  const dataLines: string[] = []
  for (const line of raw.split('\n')) {
    if (line.startsWith('event:')) event = line.slice(6).trim()
    else if (line.startsWith('data:')) dataLines.push(line.slice(5).replace(/^ /, ''))
    // 注释行（:开头）与其他字段忽略
  }
  return { event, data: dataLines.join('\n') }
}

/** 消费一个 SSE 响应，逐帧产出。调用方负责 res.ok 检查与 abort。 */
export async function* iterateSse(res: Response): AsyncGenerator<SseFrame> {
  if (!res.body) throw new Error('响应无可读流')
  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  try {
    for (;;) {
      const { done, value } = await reader.read()
      if (done) break
      buffer += decoder.decode(value, { stream: true }).replaceAll('\r\n', '\n')
      let idx: number
      while ((idx = buffer.indexOf('\n\n')) !== -1) {
        const raw = buffer.slice(0, idx)
        buffer = buffer.slice(idx + 2)
        if (raw.trim() !== '') yield parseFrame(raw)
      }
    }
    // 流结束时残留的不完整帧丢弃（正常协议下不会出现）
  } finally {
    reader.releaseLock()
  }
}
