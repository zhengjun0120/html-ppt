// 临时静态服务器：仅供本地预览 design/ 下的 demo 页面，用完即弃
const http = require('http');
const fs = require('fs');
const path = require('path');
const root = 'D:/go_files/html-ppt';
const types = { '.html': 'text/html; charset=utf-8', '.css': 'text/css', '.js': 'text/javascript', '.png': 'image/png', '.svg': 'image/svg+xml' };
http.createServer((req, res) => {
  let p = decodeURIComponent(req.url.split('?')[0]);
  if (p === '/') p = '/design/theme-demo.html';
  const f = path.join(root, p);
  fs.readFile(f, (e, d) => {
    if (e) { res.writeHead(404); res.end('not found'); return; }
    res.writeHead(200, { 'Content-Type': types[path.extname(f)] || 'application/octet-stream' });
    res.end(d);
  });
}).listen(8791, '127.0.0.1', () => console.log('serving on http://127.0.0.1:8791'));
