/*! deck-v2 editor.js — PPT 式手动编辑器（第一方文件，非 vendor，无需 PATCHES.md 登记）
 *
 * 激活：后端在 ?edit=1 的响应里注入本脚本（backend/internal/handler/inject.go）。
 * 宿主：runtime.js 的预览模式（src 带 ?preview=1&edit=1）——单页显示、preview-goto
 *       无刷新翻页、预览分支零键盘监听（runtime.js:210 提前 return），与编辑器的
 *       键盘体系零冲突。
 * 通信：postMessage，协议见 docs/deck-editor-plan.md §4.5。
 *       CSP connect-src 'none' 禁止 iframe 内发请求，保存数据只能走这条通道。
 *
 * 能力（首版，§1 决策 1）：文本自由编辑 + 拖动/缩放卡片（智能参考线 + 原位占位块）。
 * 设计核心（§1 决策 5/6）：PPT 心理模型——元素拖走即转绝对定位（append 到 slide
 * 末位，后画在上，不写 z-index），原位置留等大占位块 data-ed-placeholder，
 * 其余元素纹丝不动；占位块是真实内容，随 serialize 保留。
 */
(function () {
  'use strict';

  /* ===================== 常量与状态 ===================== */

  var SNAP_THRESHOLD = 6;   // 参考线吸附阈值（slide-local px）
  var MIN_W = 48, MIN_H = 24;
  var UNDO_CAP = 100;
  var NUDGE_BURST_MS = 600; // 方向键连续微移共享一条撤销记录

  var slides = [];
  var sel = null;            // { el, slide } 当前选中
  var hoverEl = null;
  var frameEl = null;        // 选中框 + 手柄覆盖层（挂在 slide 上）
  var editingEl = null;      // 文本编辑中的元素
  var drag = null;           // 拖动/缩放会话
  var dirty = false;
  var undoStack = [];        // { slide, html } —— 变更前该页的 innerHTML
  var redoStack = [];
  var nudgeTimer = null;

  /* ===================== 基础设施 ===================== */

  function post(msg) {
    try { window.parent.postMessage(msg, '*'); } catch (e) { /* 父页已关闭等，忽略 */ }
  }

  function setDirty(v) {
    if (dirty === v) return;
    dirty = v;
    post({ type: 'editor-dirty', dirty: dirty });
  }

  function activeSlide() {
    for (var i = 0; i < slides.length; i++) {
      if (slides[i].classList.contains('is-active')) return slides[i];
    }
    return slides[0] || null;
  }

  // .deck 被 runtime 以 --deck-scale 整体缩放：getBoundingClientRect 是屏幕像素，
  // clientWidth 是布局像素。编辑器所有几何量统一用 slide-local 布局像素。
  function slideScale(slide) {
    var r = slide.getBoundingClientRect();
    return (r.width / slide.clientWidth) || 1;
  }

  function slideLocalRect(slide, el) {
    var s = slide.getBoundingClientRect();
    var r = el.getBoundingClientRect();
    var k = slideScale(slide);
    return {
      x: (r.left - s.left) / k,
      y: (r.top - s.top) / k,
      w: r.width / k,
      h: r.height / k,
    };
  }

  /* ===================== 撤销 / 重做（按页快照 innerHTML） ===================== */
  // 全文档快照会替换 <body> 子树，把 runtime 预览模式持有的 .slide 节点引用变成
  // 游离节点（翻页立即失灵）；按页替换 innerHTML 时 section 本体不变，runtime 无感。

  function pushUndo(slide) {
    undoStack.push({ slide: slide, html: slide.innerHTML });
    if (undoStack.length > UNDO_CAP) undoStack.shift();
    redoStack.length = 0;
  }

  // innerHTML 快照可能带着当时的编辑器痕迹（拖动开始时选中框必然在场），
  // 恢复后必须扫掉，否则出现"复活的孤儿选中框"——hideFrame 只删闭包引用，
  // 删不到这些随快照回来的节点。
  function cleanSlideArtifacts(slide) {
    slide.querySelectorAll('[data-ed-frame],[data-ed-guide]').forEach(function (n) {
      n.parentNode.removeChild(n);
    });
    slide.querySelectorAll('.ed-selected,.ed-hoverable,.ed-editing').forEach(function (n) {
      n.classList.remove('ed-selected', 'ed-hoverable', 'ed-editing');
    });
    slide.querySelectorAll('[contenteditable]').forEach(function (n) {
      n.removeAttribute('contenteditable');
    });
    slide.querySelectorAll('[spellcheck="false"]').forEach(function (n) {
      n.removeAttribute('spellcheck');
    });
  }

  function undo() {
    exitTextEdit(true);
    var entry = undoStack.pop();
    if (!entry) return;
    deselect(); // 先清痕迹再拍对侧快照，否则 redo 快照里混入选中框
    redoStack.push({ slide: entry.slide, html: entry.slide.innerHTML });
    entry.slide.innerHTML = entry.html;
    cleanSlideArtifacts(entry.slide);
    invalidateCandidates();
    setDirty(true);
  }

  function redo() {
    exitTextEdit(true);
    var entry = redoStack.pop();
    if (!entry) return;
    deselect();
    undoStack.push({ slide: entry.slide, html: entry.slide.innerHTML });
    entry.slide.innerHTML = entry.html;
    cleanSlideArtifacts(entry.slide);
    invalidateCandidates();
    setDirty(true);
  }

  /* ===================== 候选元素（两层，§1 决策 6） ===================== */

  function isExcluded(el) {
    return el.tagName === 'SCRIPT' || el.tagName === 'STYLE' ||
      el.classList.contains('notes') || el.classList.contains('speaker-notes') ||
      el.classList.contains('dot') || el.classList.contains('slide-number') ||
      el.hasAttribute('data-ed-frame') || el.hasAttribute('data-ed-guide') ||
      el.hasAttribute('data-ed-placeholder') ||
      (!el.textContent.trim() && !el.querySelector('img,svg,video,canvas'));
  }

  // 候选判定（§1 决策 6 · 方案 A）：块级元素一律可选；行内元素须有盒子样式
  // （背景/边框/内边距，如 span.tag 小卡片）才可选——纯文字行内（渐变字、mono
  // 高亮）不选，想改它的字双击进编辑即可。深度不限，点击选中最深候选。
  function isCandidate(el) {
    // 讲稿区、编辑器自身的框/手柄/参考线/占位块及其子孙，永不参与候选
    if (el.closest('[data-ed-frame],[data-ed-guide],[data-ed-placeholder],.notes,.speaker-notes')) return false;
    if (isExcluded(el)) return false;
    var s = getComputedStyle(el);
    var d = s.display;
    if (d === 'none' || d === 'contents') return false;
    if (d === 'inline') {
      var bg = s.backgroundColor;
      return s.borderWidth !== '0px' ||
        (bg !== 'rgba(0, 0, 0, 0)' && bg !== 'transparent') ||
        parseFloat(s.paddingLeft) > 0 || parseFloat(s.paddingTop) > 0;
    }
    return true;
  }

  // 候选集按页缓存：hover 每次 pointermove 都要查，150+ 元素逐个 getComputedStyle
  // 不可接受。缓存随一切结构变更失效（pointerup / 退出编辑 / 撤销重做 / 转绝对定位）。
  var candCache = null; // { slide, set }

  function invalidateCandidates() {
    candCache = null;
  }

  function collectCandidates(slide) {
    if (candCache && candCache.slide === slide) return candCache.set;
    var set = new Set();
    var all = slide.querySelectorAll('*');
    for (var i = 0; i < all.length; i++) {
      if (isCandidate(all[i])) set.add(all[i]);
    }
    candCache = { slide: slide, set: set };
    return set;
  }

  function candidateFrom(target, slide) {
    var set = collectCandidates(slide);
    var el = target;
    while (el && el !== slide) {
      if (set.has(el)) return el;
      el = el.parentElement;
    }
    return null;
  }

  /* ===================== 选中框与手柄 ===================== */

  var HANDLE_DIRS = ['nw', 'n', 'ne', 'e', 'se', 's', 'sw', 'w'];

  function showFrame(el) {
    hideFrame();
    var slide = sel.slide;
    var r = slideLocalRect(slide, el);
    frameEl = document.createElement('div');
    frameEl.setAttribute('data-ed-frame', '1');
    frameEl.style.cssText = 'position:absolute;pointer-events:none;z-index:2147483646;' +
      'left:' + r.x + 'px;top:' + r.y + 'px;width:' + r.w + 'px;height:' + r.h + 'px;' +
      'outline:2px solid #6366f1;outline-offset:0;cursor:move';
    HANDLE_DIRS.forEach(function (dir) {
      var h = document.createElement('div');
      h.setAttribute('data-ed-h', dir);
      h.style.cssText = 'position:absolute;width:12px;height:12px;background:#fff;' +
        'border:2px solid #6366f1;border-radius:2px;pointer-events:auto;z-index:2147483647;' +
        'box-sizing:border-box;' + handlePos(dir);
      frameEl.appendChild(h);
    });
    slide.appendChild(frameEl);
  }

  function handlePos(dir) {
    switch (dir) {
      case 'nw': return 'left:-6px;top:-6px;cursor:nwse-resize';
      case 'n':  return 'left:calc(50% - 6px);top:-6px;cursor:ns-resize';
      case 'ne': return 'right:-6px;top:-6px;cursor:nesw-resize';
      case 'e':  return 'right:-6px;top:calc(50% - 6px);cursor:ew-resize';
      case 'se': return 'right:-6px;bottom:-6px;cursor:nwse-resize';
      case 's':  return 'left:calc(50% - 6px);bottom:-6px;cursor:ns-resize';
      case 'sw': return 'left:-6px;bottom:-6px;cursor:nesw-resize';
      case 'w':  return 'left:-6px;top:calc(50% - 6px);cursor:ew-resize';
    }
    return '';
  }

  function updateFrame() {
    if (!frameEl || !sel) return;
    var r = slideLocalRect(sel.slide, sel.el);
    frameEl.style.left = r.x + 'px';
    frameEl.style.top = r.y + 'px';
    frameEl.style.width = r.w + 'px';
    frameEl.style.height = r.h + 'px';
  }

  function hideFrame() {
    // 全局清扫：只可能有一个在册 frame，但 undo 快照恢复等历史原因可能残留
    // 游离/复制的 frame 节点，一并列掉
    document.querySelectorAll('[data-ed-frame]').forEach(function (n) {
      if (n.parentNode) n.parentNode.removeChild(n);
    });
    frameEl = null;
  }

  function select(el, slide) {
    deselect();
    sel = { el: el, slide: slide };
    el.classList.add('ed-selected');
    showFrame(el);
    postSelection();
  }

  function deselect() {
    exitTextEdit(true);
    if (hoverEl) { hoverEl.classList.remove('ed-hoverable'); hoverEl = null; }
    if (sel && sel.el) sel.el.classList.remove('ed-selected');
    sel = null;
    hideFrame();
    clearGuides();
    postSelection();
  }

  /* ===================== 选中态上报（父页工具栏据此启用字号步进器） ===================== */

  function postSelection() {
    var size = null;
    if (sel && sel.el) {
      size = Math.round(parseFloat(getComputedStyle(sel.el).fontSize)) || null;
    }
    post({ type: 'editor-selection', selected: !!sel, fontSize: size });
  }

  /* ===================== 字号缩放（§4.4：选中块及其内部文字等比缩放） ===================== */

  var FONT_MIN = 10, FONT_MAX = 200;
  var fontTimer = null; // 撤销突发合并，与方向键微移同套路

  // 选中元素自身 + 所有"直接持有文本"的后代。两层选择模型选不到卡片里的
  // 单个文字（如 grid 卡里的 h3），所以字号作用于整块——PPT 的文本块语义。
  function textLeaves(root) {
    var out = [];
    var hasDirectText = function (el) {
      for (var n = el.firstChild; n; n = n.nextSibling) {
        if (n.nodeType === 3 && n.nodeValue.trim()) return true;
      }
      return false;
    };
    if (hasDirectText(root)) out.push(root);
    var all = root.querySelectorAll('*');
    for (var i = 0; i < all.length; i++) {
      if (hasDirectText(all[i])) out.push(all[i]);
    }
    return out;
  }

  function applyFontFactor(factor) {
    if (!sel) return;
    var first = fontTimer === null;
    if (first) pushUndo(sel.slide);
    clearTimeout(fontTimer);
    fontTimer = setTimeout(function () { fontTimer = null; }, NUDGE_BURST_MS);
    var targets = textLeaves(sel.el);
    for (var i = 0; i < targets.length; i++) {
      var el = targets[i];
      var cs = getComputedStyle(el);
      var size = parseFloat(cs.fontSize) || 16;
      var next = Math.min(FONT_MAX, Math.max(FONT_MIN, Math.round(size * factor * 10) / 10));
      el.style.fontSize = next + 'px';
      // px 行高随动；无量纲/normal 行高天然等比，不动
      var lh = parseFloat(cs.lineHeight);
      if (!isNaN(lh)) el.style.lineHeight = Math.round(lh * factor * 10) / 10 + 'px';
    }
    // 盒子高度随动：转换时锁定的高度不跟着字号长，文字就会溢出框外
    if (sel.el.style.height) {
      var h = parseFloat(sel.el.style.height);
      if (!isNaN(h)) sel.el.style.height = Math.max(MIN_H, Math.round(h * factor * 10) / 10) + 'px';
    }
    healHeight(sel.el);
    updateFrame();
    setDirty(true);
    postSelection();
  }

  // 溢出自愈：锁高盒子的内容长高了（字号/改字/宽度变窄换行）就把高度贴到内容。
  // 只对纯流式内容生效——子元素里有绝对定位的（设计型盒子）内容量不可信，不动。
  function healHeight(el) {
    if (!el.style.height) return;
    var kids = el.children;
    for (var i = 0; i < kids.length; i++) {
      if (getComputedStyle(kids[i]).position === 'absolute') return;
    }
    if (el.scrollHeight > el.clientHeight + 1) {
      el.style.height = el.scrollHeight + 'px';
    }
  }

  /* ===================== 流式 → 绝对定位 + 占位块（决策 5 的核心） ===================== */

  // 元素第一次被拖/缩/微移时调用：原 DOM 位次插等大占位块顶住流式空间，
  // 元素转绝对定位（包含块 = .slide 的 padding box，slideLocalRect 直接可用），
  // append 到 slide 末位——后画的在上层，无需 z-index。
  function transformToAbsolute(el, slide) {
    if (el.style.position === 'absolute' && el.parentElement === slide) return;
    var r = slideLocalRect(slide, el);
    var m = getComputedStyle(el);

    // 占位块 = 原元素的浅克隆（类名/结构属性全保留）：布局算法看到的是和原元素
    // 完全相同的盒子（display/flex 分配/min-width/伪元素钩子类一应俱全），
    // 兄弟元素在数学上必然冻结。合成 div 只拷宽高 margin，弹性/grid 重排
    // 会挤动没被选中的兄弟（实测事故）。
    var ph = el.cloneNode(false);
    ph.removeAttribute('id');
    ph.setAttribute('data-ed-placeholder', '1');
    ph.classList.remove('ed-selected', 'ed-hoverable', 'ed-editing');
    // 克隆没有内容，盒子尺寸/弹性分配用原元素现值显式钉死
    ph.style.width = r.w + 'px';
    ph.style.height = r.h + 'px';
    ph.style.margin = m.marginTop + ' ' + m.marginRight + ' ' + m.marginBottom + ' ' + m.marginLeft;
    ph.style.flex = m.flex;
    ph.style.alignSelf = m.alignSelf;
    ph.style.boxSizing = m.boxSizing;
    // 空洞视觉中性化：皮肤透明（不影响布局），伪元素装饰关掉
    ph.style.backgroundColor = 'transparent';
    ph.style.borderColor = 'transparent';
    ph.style.boxShadow = 'none';
    ph.style.color = 'transparent';
    // grid 手工定位项原样带走（罕见，但带走无害）
    if (el.style.gridColumn) ph.style.gridColumn = el.style.gridColumn;
    if (el.style.gridRow) ph.style.gridRow = el.style.gridRow;
    el.parentNode.insertBefore(ph, el);

    el.style.position = 'absolute';
    el.style.left = r.x + 'px';
    el.style.top = r.y + 'px';
    el.style.width = r.w + 'px';
    el.style.height = r.h + 'px';
    el.style.margin = '0';
    slide.appendChild(el);
    invalidateCandidates(); // DOM 结构变了（占位块+搬家），候选集重建
  }

  /* ===================== 智能参考线 ===================== */

  var guideLayer = null;

  function clearGuides() {
    if (guideLayer) {
      while (guideLayer.firstChild) guideLayer.removeChild(guideLayer.firstChild);
    }
  }

  function ensureGuideLayer(slide) {
    if (!guideLayer || guideLayer.parentElement !== slide) {
      if (guideLayer && guideLayer.parentElement) guideLayer.parentElement.removeChild(guideLayer);
      guideLayer = document.createElement('div');
      guideLayer.setAttribute('data-ed-guide', 'layer');
      guideLayer.style.cssText = 'position:absolute;inset:0;pointer-events:none;z-index:2147483645';
      slide.appendChild(guideLayer);
    }
    return guideLayer;
  }

  function drawGuide(slide, dir, pos) {
    var layer = ensureGuideLayer(slide);
    var g = document.createElement('div');
    g.setAttribute('data-ed-guide', 'line');
    g.style.cssText = 'position:absolute;pointer-events:none;background:#f59e0b;' +
      (dir === 'v'
        ? 'left:' + pos + 'px;top:0;width:1px;height:100%'
        : 'top:' + pos + 'px;left:0;height:1px;width:100%');
    layer.appendChild(g);
  }

  // 参照集：同页其他候选元素的三线（首/中/尾）+ 内容区边（padding 72/96）+ 画布中线
  function snapRects(slide, exclude) {
    var refs = { xs: [], ys: [] };
    var sw = slide.clientWidth, sh = slide.clientHeight;
    var addX = function (v) { refs.xs.push(v); };
    var addY = function (v) { refs.ys.push(v); };
    var set = collectCandidates(slide);
    set.forEach(function (el) {
      if (el === exclude) return;
      var r = slideLocalRect(slide, el);
      addX(r.x); addX(r.x + r.w / 2); addX(r.x + r.w);
      addY(r.y); addY(r.y + r.h / 2); addY(r.y + r.h);
    });
    addX(96); addX(sw - 96); addX(sw / 2);   // .slide padding:72px 96px
    addY(72); addY(sh - 72); addY(sh / 2);
    return refs;
  }

  // 对 {lead 数组} vs {refs} 找最近吸附对，命中返回 {value, guide}
  function snapAxis(leads, refs) {
    var best = null;
    leads.forEach(function (p) {
      refs.forEach(function (rv) {
        var d = Math.abs(p - rv);
        if (d <= SNAP_THRESHOLD && (!best || d < best.d)) best = { d: d, value: rv, at: p };
      });
    });
    return best;
  }

  /* ===================== 拖动 / 缩放 ===================== */

  // 角手柄缩放的字号基准：选中元素内所有"直接持有文本"的节点的现值
  function fontBasesOf(root) {
    return textLeaves(root).map(function (el) {
      var cs = getComputedStyle(el);
      return { el: el, size: parseFloat(cs.fontSize) || 16, lh: parseFloat(cs.lineHeight) };
    });
  }

  function beginDrag(e, el, slide) {
    var k = slideScale(slide);
    var r = slideLocalRect(slide, el);
    drag = {
      kind: 'drag', el: el, slide: slide, k: k,
      start: { x: r.x, y: r.y, w: r.w, h: r.h },
      px: e.clientX, py: e.clientY,
      moved: false, refs: snapRects(slide, el),
    };
  }

  function beginResize(e, dir, el, slide) {
    // 与拖动同一条纪律：先转绝对定位再改尺寸。流式元素上直接写 width/height
    // 会引发居中布局（.slide 是 justify-content:center）整页回流——元素带着
    // 选中框一起"跳走"；且 left/top 对流式元素不生效，w/n 手柄的数学全落空。
    transformToAbsolute(el, slide);
    var k = slideScale(slide);
    var r = slideLocalRect(slide, el);
    drag = {
      kind: 'resize', dir: dir, el: el, slide: slide, k: k,
      start: { x: r.x, y: r.y, w: r.w, h: r.h },
      px: e.clientX, py: e.clientY,
      moved: false, refs: null,
      // 角手柄 = 内容等比缩放（像缩放图片）：以起始字号为基准、按宽度比例实时换算
      fontBases: dir.length === 2 ? fontBasesOf(el) : null,
    };
  }

  function onPointerMove(e) {
    if (drag) {
      var dx = (e.clientX - drag.px) / drag.k;
      var dy = (e.clientY - drag.py) / drag.k;
      if (!drag.moved && Math.abs(dx) < 0.5 && Math.abs(dy) < 0.5) return;
      if (!drag.moved) {
        drag.moved = true;
        pushUndo(drag.slide);
        document.body.classList.add('ed-busy');
        if (drag.kind === 'drag') transformToAbsolute(drag.el, drag.slide);
      }

      var s = drag.start;
      if (drag.kind === 'drag') {
        var nx = s.x + dx, ny = s.y + dy;
        var sw = drag.slide.clientWidth, sh = drag.slide.clientHeight;
        // 智能参考线：仅拖动时，先吸附后钳位（贴边优先）
        clearGuides();
        var sx = snapAxis([nx, nx + s.w / 2, nx + s.w], drag.refs.xs);
        var sy = snapAxis([ny, ny + s.h / 2, ny + s.h], drag.refs.ys);
        if (sx) { nx += sx.value - sx.at; drawGuide(drag.slide, 'v', sx.value); }
        if (sy) { ny += sy.value - sy.at; drawGuide(drag.slide, 'h', sy.value); }
        nx = Math.min(Math.max(nx, 0), Math.max(0, sw - s.w));
        ny = Math.min(Math.max(ny, 0), Math.max(0, sh - s.h));
        drag.el.style.left = nx + 'px';
        drag.el.style.top = ny + 'px';
      } else {
        var r = { x: s.x, y: s.y, w: s.w, h: s.h };
        if (drag.dir.indexOf('e') >= 0) r.w = Math.max(MIN_W, s.w + dx);
        if (drag.dir.indexOf('s') >= 0) r.h = Math.max(MIN_H, s.h + dy);
        if (drag.dir.indexOf('w') >= 0) {
          r.w = Math.max(MIN_W, s.w - dx);
          r.x = s.x + (s.w - r.w);
        }
        if (drag.dir.indexOf('n') >= 0) {
          r.h = Math.max(MIN_H, s.h - dy);
          r.y = s.y + (s.h - r.h);
        }
        drag.el.style.left = r.x + 'px';
        drag.el.style.top = r.y + 'px';
        drag.el.style.width = r.w + 'px';
        drag.el.style.height = r.h + 'px';
        // 角手柄：文字随盒子等比缩放（基准是拖动开始时的现值，避免连乘漂移）
        if (drag.fontBases && s.w > 0) {
          var ff = r.w / s.w;
          if (isFinite(ff) && ff > 0) {
            for (var bi = 0; bi < drag.fontBases.length; bi++) {
              var fb = drag.fontBases[bi];
              fb.el.style.fontSize =
                Math.min(FONT_MAX, Math.max(FONT_MIN, Math.round(fb.size * ff * 10) / 10)) + 'px';
              if (!isNaN(fb.lh)) fb.el.style.lineHeight = Math.round(fb.lh * ff * 10) / 10 + 'px';
            }
            postSelection();
          }
        }
      }
      updateFrame();
      return;
    }

    // 非拖动态：hover 提示（文本编辑中光标归 caret，不刷 outline）
    if (editingEl) return;
    var slide = activeSlide();
    if (!slide) return;
    var c = candidateFrom(e.target, slide);
    if (c !== hoverEl) {
      if (hoverEl) hoverEl.classList.remove('ed-hoverable');
      hoverEl = c && (!sel || c !== sel.el) ? c : null;
      if (hoverEl) hoverEl.classList.add('ed-hoverable');
    }
  }

  function onPointerUp() {
    if (!drag) return;
    document.body.classList.remove('ed-busy');
    if (drag.moved) {
      clearGuides();
      // 宽度变了的缩放（含角手柄）改变了换行：松手后把盒子高度贴回内容。
      // 纵向手柄（n/s）是用户在控高，不治——拖小弹回会跟人打架。
      if (drag.kind === 'resize' && (drag.dir.indexOf('e') >= 0 || drag.dir.indexOf('w') >= 0)) {
        healHeight(drag.el);
        updateFrame();
      }
      invalidateCandidates();
      setDirty(true);
    }
    drag = null;
  }

  function onPointerDown(e) {
    if (e.button !== 0) return;
    var slide = activeSlide();
    if (!slide) return;

    // 手柄优先：resize 不走候选命中
    var handle = e.target.closest ? e.target.closest('[data-ed-h]') : null;
    if (handle && frameEl) {
      e.preventDefault();
      beginResize(e, handle.getAttribute('data-ed-h'), sel.el, sel.slide);
      return;
    }
    if (editingEl) return; // 文本编辑中，pointerdown 归 caret

    var c = candidateFrom(e.target, slide);
    if (c) {
      e.preventDefault();
      if (sel && sel.el === c) {
        // 已选中元素上按下：直接进入拖动
        beginDrag(e, c, slide);
      } else {
        select(c, slide);
      }
    } else if (!frameEl || !frameEl.contains(e.target)) {
      deselect();
    }
  }

  /* ===================== 文本编辑 ===================== */

  function enterTextEdit(el) {
    if (editingEl === el) return;
    exitTextEdit(true);
    pushUndo(el.closest('.slide') || activeSlide());
    editingEl = el;
    el.classList.remove('ed-hoverable');
    el.classList.add('ed-editing');
    // contentEditable='true'（而非 plaintext-only）：后者会把 Enter/脚本插的
    // <br> 都吞成纯文本 \n。富文本风险由 paste 拦截 + Ctrl+B/I/U 拦截兜住。
    el.contentEditable = 'true';
    // 行内 normal 压过 UA 给编辑态的 pre-wrap：源码结构性换行在编辑态也不可见，
    // 编辑态渲染 = 最终渲染（行内样式必赢 UA 表，class 选择器压不住它）
    el.style.whiteSpace = 'normal';
    el.setAttribute('spellcheck', 'false');
    el.focus();
  }

  // plaintext-only 编辑态默认 pre-wrap：源码里标签之间的结构性换行会显示成
  // 空行，退出时再被归一化成 <br> 就成了凭空多出的空行（实测事故）。
  // 正解是让编辑态渲染 = 最终渲染：white-space:normal 让结构性 \n 从头到尾
  // 不可见；用户的 Enter 由 Enter 处理器手动插真 <br>，粘贴多行拆行插 <br>。
  function insertBrAtCaret() {
    var s = window.getSelection();
    if (!s || !s.rangeCount) return;
    var r = s.getRangeAt(0);
    r.deleteContents();
    var br = document.createElement('br');
    r.insertNode(br);
    r.setStartAfter(br);
    r.collapse(true);
    s.removeAllRanges();
    s.addRange(r);
  }

  // 粘贴多行纯文本：拆行插真 <br>（编辑态 normal 折叠 \n，必须转成 br 才可见）
  function onPaste(e) {
    if (!editingEl) return;
    var clip = e.clipboardData || window.clipboardData;
    if (!clip) return;
    var txt = clip.getData('text/plain');
    if (txt == null) return;
    e.preventDefault();
    var lines = txt.replace(/\r/g, '').split('\n');
    for (var i = 0; i < lines.length; i++) {
      if (i > 0) insertBrAtCaret();
      if (lines[i]) document.execCommand('insertText', false, lines[i]);
    }
  }

  function exitTextEdit(silent) {
    if (!editingEl) return;
    var el = editingEl;
    editingEl = null;
    el.removeAttribute('contenteditable');
    el.removeAttribute('spellcheck');
    el.classList.remove('ed-editing');
    el.style.removeProperty('white-space'); // 还原模板自己的 white-space 规则
    try { el.blur(); } catch (e) { /* 已失焦 */ }
    healHeight(el); // 打字换行后盒子高度贴回内容（所见即所得）
    updateFrame();
    invalidateCandidates(); // 编辑可能插了 <br>、heal 改了尺寸
    if (silent !== true) {
      setDirty(true);
    }
  }

  function onDblClick(e) {
    var slide = activeSlide();
    if (!slide) return;
    var c = candidateFrom(e.target, slide);
    if (!c) return;
    if (!sel || sel.el !== c) select(c, slide);
    enterTextEdit(c);
  }

  /* ===================== 键盘 ===================== */

  function nudge(dx, dy) {
    if (!sel) return;
    var first = nudgeTimer === null;
    if (first) pushUndo(sel.slide);
    clearTimeout(nudgeTimer);
    nudgeTimer = setTimeout(function () { nudgeTimer = null; }, NUDGE_BURST_MS);
    if (first) transformToAbsolute(sel.el, sel.slide);
    var x = parseFloat(sel.el.style.left) || 0;
    var y = parseFloat(sel.el.style.top) || 0;
    sel.el.style.left = (x + dx) + 'px';
    sel.el.style.top = (y + dy) + 'px';
    updateFrame();
    setDirty(true);
  }

  function onKeyDown(e) {
    var meta = e.ctrlKey || e.metaKey;

    if (meta && (e.key === 'z' || e.key === 'Z')) {
      if (editingEl) return; // 编辑中放行浏览器原生 undo
      e.preventDefault();
      if (e.shiftKey) redo(); else undo();
      return;
    }
    if (meta && (e.key === 'y' || e.key === 'Y')) {
      if (editingEl) return;
      e.preventDefault();
      redo();
      return;
    }
    if (meta && (e.key === 's' || e.key === 'S')) {
      e.preventDefault();
      exitTextEdit(true);
      post({ type: 'editor-serialize', html: serialize() });
      return;
    }
    if (meta && (e.key === '=' || e.key === '+')) {
      if (!sel || editingEl) return;
      e.preventDefault();
      applyFontFactor(1.1);
      return;
    }
    if (meta && e.key === '-') {
      if (!sel || editingEl) return;
      e.preventDefault();
      applyFontFactor(0.9);
      return;
    }
    if (e.key === 'Escape') {
      if (editingEl) { exitTextEdit(); return; }
      if (sel && sel.el && collectCandidates(sel.slide).has(sel.el)) {
        var set = collectCandidates(sel.slide);
        var parentEl = sel.el.parentElement;
        // Esc 链：level2 → 上选其父（level1）；level1 → 取消选中
        if (parentEl && parentEl !== sel.slide && set.has(parentEl)) {
          select(parentEl, sel.slide);
        } else {
          deselect();
        }
      }
      return;
    }
    if (editingEl) {
      if (e.key === 'Enter') {
        // 'true' 模式下 execCommand 产真 <br>；'plaintext-only' 会吞成 \n（已弃用）
        e.preventDefault();
        document.execCommand('insertLineBreak');
      } else if (meta && (e.key === 'b' || e.key === 'i' || e.key === 'u' || e.key === 'B' || e.key === 'I' || e.key === 'U')) {
        // 防富文本混入：加粗/斜体/下划线一律不放行
        e.preventDefault();
      }
      return;
    }
    if (!sel) return;
    var step = e.shiftKey ? 10 : 1;
    switch (e.key) {
      case 'ArrowLeft': e.preventDefault(); nudge(-step, 0); break;
      case 'ArrowRight': e.preventDefault(); nudge(step, 0); break;
      case 'ArrowUp': e.preventDefault(); nudge(0, -step); break;
      case 'ArrowDown': e.preventDefault(); nudge(0, step); break;
    }
  }

  /* ===================== serialize（在克隆上清理，活动 DOM 不动） ===================== */

  function serialize() {
    // 先退出编辑态：plaintext-only 的 \n 换行必须归一化成 <br> 才能进文件，
    // 否则"敲着换行直接点保存"会把换行丢掉。healHeight 同理顺带生效。
    exitTextEdit(true);
    var root = document.documentElement.cloneNode(true);

    // 1. 编辑器自身痕迹
    // 双保险用 querySelectorAll：历史文件里可能已经落盘过一个 #ed-styles
    // （v000002 双实例事故），加上本次注入的，同屏可能有两份
    root.querySelectorAll('#ed-styles').forEach(function (n) {
      n.parentNode.removeChild(n);
    });
    // 注入的 editor.js <script> 是"每请求由服务端注入"的运行态，不是内容——
    // 不剥掉就会随保存落盘，下次 ?edit=1 加载被双重注入（所有事件跑两遍）
    root.querySelectorAll('script[src="/assets/deck-v2/editor.js"]').forEach(function (n) {
      n.parentNode.removeChild(n);
    });
    root.querySelectorAll('[data-ed-frame],[data-ed-guide]').forEach(function (n) {
      n.parentNode.removeChild(n);
    });
    root.querySelectorAll('.ed-selected,.ed-hoverable,.ed-editing').forEach(function (n) {
      n.classList.remove('ed-selected', 'ed-hoverable', 'ed-editing');
    });
    root.querySelectorAll('[contenteditable]').forEach(function (n) {
      n.removeAttribute('contenteditable');
    });
    root.querySelectorAll('[spellcheck="false"]').forEach(function (n) {
      n.removeAttribute('spellcheck');
    });

    // 2. runtime 运行态 → 还原到"生成产物"的原始形状
    root.querySelectorAll('.slide').forEach(function (s) {
      s.classList.remove('is-active', 'is-prev');
    });
    root.removeAttribute('data-preview');
    root.removeAttribute('style'); // runtime syncAmbient 同时给 <html> 和 <body> 写背景色
    var body = root.querySelector('body');
    if (body) {
      body.removeAttribute('data-preview');
      body.removeAttribute('style'); // runtime syncAmbient 写的 background-color 等运行态
      body.classList.remove('ed-busy');
    }
    var deck = root.querySelector('.deck');
    if (deck) deck.removeAttribute('style'); // --deck-scale/--deck-w/--deck-h/--ambient-*
    root.querySelectorAll('.notes,.speaker-notes').forEach(function (n) {
      n.removeAttribute('style'); // 预览分支打的 inline display:none（base.css 已有 !important 规则）
    });
    // data-ed-placeholder 保留！占位块是编辑结果的一部分（§1 决策 5）。
    // data-current/data-total 保留（runtime 重算，无害）。

    return '<!DOCTYPE html>\n' + root.outerHTML;
  }

  /* ===================== 消息协议（§4.5） ===================== */

  var MSG_FROM_PARENT = { 'editor-save': 1, 'editor-saved': 1, 'editor-undo': 1, 'editor-redo': 1, 'editor-font': 1 };

  function onMessage(e) {
    if (e.source !== window.parent) return;
    var d = e.data;
    if (!d || typeof d !== 'object' || !MSG_FROM_PARENT[d.type]) return;
    switch (d.type) {
      case 'editor-save':
        post({ type: 'editor-serialize', html: serialize() });
        break;
      case 'editor-saved':
        setDirty(false);
        break;
      case 'editor-undo': undo(); break;
      case 'editor-redo': redo(); break;
      case 'editor-font':
        if (typeof d.factor === 'number' && d.factor > 0) applyFontFactor(d.factor);
        break;
    }
  }

  /* ===================== 样式注入 ===================== */

  function injectStyles() {
    // 自愈：历史保存产物里可能残留一份 #ed-styles（双实例事故的产物），
    // 先清掉再注入，保证全文档只有一份
    document.querySelectorAll('#ed-styles').forEach(function (n) {
      n.parentNode.removeChild(n);
    });
    var style = document.createElement('style');
    style.id = 'ed-styles';
    style.textContent =
      '.ed-selected{outline:2px solid #6366f1;outline-offset:2px}' +
      '.ed-hoverable{outline:1px dashed rgba(99,102,241,.5)}' +
      '.ed-editing{outline:2px dashed #6366f1;cursor:text;white-space:normal}' +
      'body.ed-busy,body.ed-busy *{cursor:move!important;user-select:none!important;-webkit-user-select:none!important}' +
      '[data-ed-placeholder]::before,[data-ed-placeholder]::after{content:none!important}';
    document.head.appendChild(style);
  }

  /* ===================== 启动 ===================== */

  function init() {
    slides = Array.prototype.slice.call(document.querySelectorAll('.deck > section.slide'));
    if (!slides.length) slides = Array.prototype.slice.call(document.querySelectorAll('section.slide'));
    if (!slides.length) return; // 不是 deck 页面，安静退出

    injectStyles();
    document.addEventListener('pointerdown', onPointerDown, true);
    document.addEventListener('pointermove', onPointerMove);
    document.addEventListener('pointerup', onPointerUp);
    document.addEventListener('dblclick', onDblClick);
    document.addEventListener('keydown', onKeyDown, true);
    document.addEventListener('paste', onPaste, true);
    window.addEventListener('message', onMessage);
    window.addEventListener('beforeunload', function () {
      exitTextEdit(true);
    });

    post({ type: 'editor-ready', pages: slides.length });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
