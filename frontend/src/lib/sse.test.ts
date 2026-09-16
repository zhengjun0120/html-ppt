import { describe, expect, it } from 'vitest'

import { iterateSse } from './sse'

function sseResponse(chunks: string[]): Response {
  const encoder = new TextEncoder()
  let i = 0
  return new Response(null, {})
    ? new Response(
        new ReadableStream({
          pull(controller) {
            if (i < chunks.length) controller.enqueue(encoder.encode(chunks[i++]))
            else controller.close()
          },
        }),
      )
    : (undefined as never)
}

describe('iterateSse（POST SSE 解析器）', () => {
  it('解析单 chunk 内的多帧', async () => {
    const res = sseResponse(['event: delta\ndata: {"type":"delta"}\n\nevent: done\ndata: {"type":"done"}\n\n'])
    const frames = []
    for await (const f of iterateSse(res)) frames.push(f)
    expect(frames).toEqual([
      { event: 'delta', data: '{"type":"delta"}' },
      { event: 'done', data: '{"type":"done"}' },
    ])
  })

  it('帧跨 chunk 分片时正确缓冲拼接', async () => {
    const res = sseResponse([
      'event: del',
      'ta\ndata: {"a":',
      '1}\n\nevent: session\ndata: {"type":"session"}\n\n',
    ])
    const frames = []
    for await (const f of iterateSse(res)) frames.push(f)
    expect(frames).toEqual([
      { event: 'delta', data: '{"a":1}' },
      { event: 'session', data: '{"type":"session"}' },
    ])
  })

  it('兼容 CRLF 分隔', async () => {
    const res = sseResponse(['event: delta\r\ndata: x\r\n\r\n'])
    const frames = []
    for await (const f of iterateSse(res)) frames.push(f)
    expect(frames).toEqual([{ event: 'delta', data: 'x' }])
  })
})
