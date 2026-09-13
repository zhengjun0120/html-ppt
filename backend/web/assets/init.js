// deck 初始化脚本：所有 deck 共用这一份（第①层：框架，LLM 不可见不可改）。
// 将来交互组件（翻转卡片、拖拽排序等）的行为绑定也统一加在这里。

// ---------- 主题与翻页动画 ----------
// 翻页动画不在 CSS 变量体系里，从 deck 内的主题配置块读取；
// 配置缺失/损坏时退回默认，框架层永不因配置挂掉
let transition = 'slide';
try {
  transition = JSON.parse(document.getElementById('deck-theme').textContent).transition || transition;
} catch (e) {}

Reveal.initialize({
  hash: true,       // URL 带 #/2 页码：刷新不丢位置，截图也能精确定位到某页
  transition: transition,
  controls: true,
  progress: true,
});

// ---------- 溢出兜底（fit）----------
// 为什么写在 init.js 而不是单独的 fit.js：骨架只在 write_deck 时固化一次，
// 新增 <script> 标签对已经存在的 deck 无效；init.js 是所有 deck 都已经在引用的那一份。
//
// 策略两级，**都不裁切内容**（宁可变小，也不用 overflow:hidden 把内容藏起来）：
//   一级：缩 section 的 font-size。组件库的尺寸与间距全是 em（--space-md: 0.6em、
//         font-size: 0.62em），所以缩字号会让文字、内边距、卡片一起等比缩小，
//         并借重排把宽度重新用满。
//   二级：缩到下限仍装不下（内容量级明显失控）时，把内容包进内层 wrapper 整体缩放。
//         不写在 section 自己身上——reveal 翻页时会改写 section 的 transform，
//         写在它上面会被覆盖并造成跳动。
(function () {
  var FIT_MIN = 0.55;      // 一级缩放的下限
  var FIT_STEP = 0.94;     // 每步收缩比例
  var FIT_TOL = 1.02;      // 2% 容差，避免浮点误差误判为溢出
  var FIT_MAX_STEPS = 16;
  var FIT_MARGIN = 32;     // 画布四周留的安全余量（px）：正好塞满 700 高会显得很挤

  function config(key) {
    try {
      var cfg = Reveal.getConfig() || {};
      return (key === 'h' ? cfg.height : cfg.width) || (key === 'h' ? 700 : 960);
    } catch (e) {
      return key === 'h' ? 700 : 960;
    }
  }

  // reveal 的 section 是"内容撑开"的（没有固定高度），所以溢出判定必须拿
  // 内容尺寸去比 **配置的逻辑画布**，不能用 scrollHeight/clientHeight（那俩永远相等）。
  // 非当前页是 display:none，量之前要临时显形（对用户不可见），量完还原。
  function withVisible(sec, fn) {
    var hidden = !sec.classList.contains('present');
    if (hidden) {
      sec.style.display = 'block';
      sec.style.visibility = 'hidden';
    }
    var out = fn();
    if (hidden) {
      sec.style.removeProperty('display');
      sec.style.removeProperty('visibility');
    }
    return out;
  }

  function overflowRatio(sec) {
    return withVisible(sec, function () {
      // 纵向：section 是内容撑开的，scrollHeight 就是内容高度，拿去比画布高（留余量）
      // 横向：section 宽度本来就是满画布（width:100%），所以只能拿内容比自己盒子的宽，
      //       不能减余量也不能比画布——否则每一页都会被判成横向溢出
      var availH = Math.max(1, config('h') - FIT_MARGIN);
      var wRatio = sec.scrollWidth / Math.max(1, sec.clientWidth);
      return Math.max(sec.scrollHeight / availH, wRatio);
    });
  }

  // 还原到未 fit 的状态：fit 必须可重复执行（字体加载完、内容变化后要重跑）
  function resetFit(sec) {
    sec.style.removeProperty('font-size');
    var inner = sec.querySelector(':scope > .ppt-fit-inner');
    if (inner) {
      while (inner.firstChild) sec.insertBefore(inner.firstChild, inner);
      inner.remove();
    }
  }

  function fitSection(sec) {
    resetFit(sec);
    if (overflowRatio(sec) <= FIT_TOL) return;   // 绝大多数页到此结束

    // 一级：等比缩字号
    var scale = 1;
    for (var i = 0; i < FIT_MAX_STEPS && scale > FIT_MIN; i++) {
      scale = Math.max(FIT_MIN, scale * FIT_STEP);
      sec.style.fontSize = (scale * 100).toFixed(1) + '%';
      if (overflowRatio(sec) <= FIT_TOL) return;
    }

    // 二级：整页内容缩放。transform 不改布局，系数能一次算出；宽度补偿会让文字
    // 重排（行变长、行数变少），所以要算两遍收敛。
    // 宽度补偿后 wrapper 比画布宽，缩放原点必须是左上，缩放后才正好铺满画布。
    var inner = document.createElement('div');
    inner.className = 'ppt-fit-inner';
    inner.style.transformOrigin = 'top left';
    while (sec.firstChild) inner.appendChild(sec.firstChild);
    sec.appendChild(inner);

    for (var pass = 0; pass < 2; pass++) {
      var availH = Math.max(1, config('h') - FIT_MARGIN);
      // 必须量在"临时显形"窗口内：隐藏页的 offsetHeight 是 0，会算出 k=1 而跳过缩放
      var layoutH = withVisible(sec, function () { return inner.offsetHeight; }) || 1;
      var k = Math.min(1, availH / layoutH);
      if (k >= 1) break;
      inner.style.transform = 'scale(' + k.toFixed(4) + ')';
      inner.style.width = (100 / k).toFixed(2) + '%';
    }
  }

  function fitAll() {
    var slides = [];
    try { slides = Reveal.getSlides() || []; } catch (e) { return; }
    for (var i = 0; i < slides.length; i++) {
      var sec = slides[i];
      if (!sec || sec.tagName !== 'SECTION') continue;
      try { fitSection(sec); } catch (e) { /* 单页出错不拖垮整个 deck */ }
    }
    try { Reveal.layout(); } catch (e) {}   // 尺寸变了，让 reveal 重算居中与整体缩放
  }

  try {
    if (Reveal.isReady && Reveal.isReady()) fitAll();
    else Reveal.on('ready', fitAll);
    // 翻到某页时再补一次：兜住"加载后才长高的内容"
    Reveal.on('slidechanged', function (e) {
      try { if (e && e.currentSlide) fitSection(e.currentSlide); } catch (err) {}
    });
  } catch (e) {}

  // 字体异步加载，加载前后的行高不同：就绪后重跑一遍结果才稳定
  if (document.fonts && document.fonts.ready && document.fonts.ready.then) {
    document.fonts.ready.then(function () { try { fitAll(); } catch (e) {} });
  }
})();
