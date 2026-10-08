/**
 * 定制对话的用户附图：本地压缩成模型友好的 data URL。
 *
 * 服务端的限制（backend usertpl/customize_images.go）：单张解码后 ≤4MB、
 * 一次 ≤3 张、PNG/JPEG/WebP。前端把长边压到 1568px 的 JPEG（q0.85），
 * 既稳稳落在服务端限制内，也把每张图的 token 成本控制在 ~1-2k。
 */

/** 与服务端 maxUserImages 对齐的每条消息附图上限 */
export const MAX_CHAT_IMAGES = 3

/** 压缩目标：长边像素（视觉输入的甜点区，再大模型也看不清） */
const MAX_EDGE = 1568

const JPEG_QUALITY = 0.85

export async function fileToChatImage(file: File): Promise<string> {
  if (!/^image\/(png|jpeg|webp)$/.test(file.type)) {
    throw new Error('只支持 PNG/JPEG/WebP 图片')
  }
  if (file.size > 10 * 1024 * 1024) {
    throw new Error('图片超过 10MB，请先压缩再发')
  }
  // 走 Image + object URL，不走 createImageBitmap——后者在部分内嵌 webview 上
  // 对合法 PNG 报 InvalidStateError（2026-10-06 实测），img 元素解码则哪里都活。
  // 浏览器解码 <img src=objectURL> 时会顺带应用 EXIF 方向，手机竖拍照片不横倒。
  const url = URL.createObjectURL(file)
  try {
    const img = await new Promise<HTMLImageElement>((resolve, reject) => {
      const el = new Image()
      el.onload = () => resolve(el)
      el.onerror = () => reject(new Error('图片无法解码，文件可能已损坏'))
      el.src = url
    })
    const scale = Math.min(1, MAX_EDGE / Math.max(img.naturalWidth, img.naturalHeight))
    const w = Math.max(1, Math.round(img.naturalWidth * scale))
    const h = Math.max(1, Math.round(img.naturalHeight * scale))
    const canvas = document.createElement('canvas')
    canvas.width = w
    canvas.height = h
    const ctx = canvas.getContext('2d')
    if (!ctx) throw new Error('画布不可用')
    ctx.fillStyle = '#ffffff' // JPEG 没有透明通道，白底兜底（PNG 透明区不发黑）
    ctx.fillRect(0, 0, w, h)
    ctx.drawImage(img, 0, 0, w, h)
    return canvas.toDataURL('image/jpeg', JPEG_QUALITY)
  } finally {
    URL.revokeObjectURL(url)
  }
}
