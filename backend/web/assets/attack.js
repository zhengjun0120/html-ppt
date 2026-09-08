// 沙箱攻击测试脚本。它是后端自己放的测试夹具（同源外链），用来验证：
//   1. iframe sandbox（不给 allow-same-origin）拦 DOM/存储/cookie
//   2. CSP connect-src 'none' 拦网络请求
// 预期：在 chat-test 的沙箱预览里五项全"被拦"；用"预览↗"裸开时
// 存储和 cookie 会"得手"（应用同源存储！）——这就是阶段5 要迁独立 origin 的实证。
// 测试完成后可删除本文件和 deck-0003。
const log = [];
const add = (name, ok, detail) =>
  log.push((ok ? "✅ 得手  " : "⛔ 被拦  ") + name + " — " + detail);

// 1. 父页面 DOM：沙箱应拦（不透明 origin 读不到跨源 document）
try { add("读父页面 DOM", true, "title=" + window.parent.document.title); }
catch (e) { add("读父页面 DOM", false, e.name); }

// 2. localStorage：沙箱应拦（SecurityError）
try { localStorage.setItem("attack", "1"); add("写 localStorage", true, "成功"); }
catch (e) { add("写 localStorage", false, e.name); }

// 3. cookie：沙箱应拦
try { add("读 cookie", true, document.cookie || "(空)"); }
catch (e) { add("读 cookie", false, e.name); }

// 4. 调本站 API：CSP connect-src 'none' 应拦
// 5. 向外网渗出：CSP 应拦（no-cors 模式下若没有 CSP，请求会真的发出去）
const tasks = [
  fetch("/api/decks").then(r => ["调本站 API", true, "HTTP " + r.status])
                     .catch(e => ["调本站 API", false, e.name + "（CSP）"]),
  fetch("https://httpbin.org/get?t=1", { mode: "no-cors" })
      .then(() => ["外渗请求", true, "请求已发出"])
      .catch(e => ["外渗请求", false, e.name + "（CSP）"]),
];

Promise.allSettled(tasks).then(rs => {
  rs.forEach(r => add(...r.value));
  document.getElementById("log").textContent = log.join("\n");
});
