// 一次性视觉自查工具：起静态服务 + chromedp 按页截图 layouts-test.html。
// 用法：go run ./cmd/shotdemo <outDir> ?preset=...（query 可选）
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	outDir := os.Args[1]
	query := ""
	if len(os.Args) > 2 {
		query = os.Args[2]
	}
	webRoot := "web"
	if _, err := os.Stat(webRoot); err != nil {
		webRoot = filepath.Join("..", "web") // 从 backend 目录外运行时
	}
	assets := filepath.Join(webRoot, "assets")
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assets))))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(assets, "layouts-test.html"))
		})
		log.Fatal(http.ListenAndServe("127.0.0.1:18621", mux))
	}()
	time.Sleep(400 * time.Millisecond)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancelT := context.WithTimeout(ctx, 90*time.Second)
	defer cancelT()

	var slides int
	url := "http://127.0.0.1:18621/layouts-test.html" + query
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1280, 800),
		chromedp.Navigate(url),
		chromedp.WaitReady("div.reveal"),
		chromedp.Sleep(2500*time.Millisecond),
		chromedp.Evaluate(`document.querySelectorAll('.slides section').length`, &slides),
	); err != nil {
		log.Fatalf("load: %v", err)
	}
	fmt.Println("slides:", slides)
	os.MkdirAll(outDir, 0o755)
	for i := 0; i < slides; i++ {
		var buf []byte
		err := chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(`Reveal.slide(%d)`, i), nil),
			chromedp.Sleep(700*time.Millisecond),
			chromedp.FullScreenshot(&buf, 80),
		)
		if err != nil {
			log.Printf("slide %d: %v", i, err)
			continue
		}
		name := filepath.Join(outDir, fmt.Sprintf("slide-%02d%s.png", i+1, tag(query)))
		if err := os.WriteFile(name, buf, 0o644); err != nil {
			log.Println(err)
		}
	}
	fmt.Println("done ->", outDir)
}

func tag(q string) string {
	if q == "" {
		return ""
	}
	return "-" + q[len("?preset="):]
}
