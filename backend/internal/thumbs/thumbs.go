// Package thumbs 管理 deck 的整本页级缩略图缓存。
//
// 缓存 = data/decks/<id>/thumbs/{1..N}.png + .stamp。有效性以 index.html 的
// （mtime+size）指纹为准：deck 任何一次页面写入都会改变它的 mtime，读取方
// 比对 stamp 不一致就整本重渲。失效由 deck 服务的写路径调用 Invalidate 完成，
// 两边机制互补——stamp 防漏、Invalidate 防陈旧目录堆积。
package thumbs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Dir 某份 deck 的缩略图目录。
func Dir(decksDir, deckID string) string {
	return filepath.Join(decksDir, deckID, "thumbs")
}

// stampContent 计算 index.html 的有效性指纹。
func stampContent(indexPath string) string {
	st, err := os.Stat(indexPath)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d|%d", st.ModTime().UnixNano(), st.Size())
}

// Valid 缩略图缓存是否仍然有效（stamp 与 index.html 当前指纹一致，且至少有第 1 页）。
func Valid(dir, indexPath string) bool {
	want := stampContent(indexPath)
	if want == "" {
		return false
	}
	b, err := os.ReadFile(filepath.Join(dir, ".stamp"))
	if err != nil || string(b) != want {
		return false
	}
	_, err = os.Stat(filepath.Join(dir, "1.png"))
	return err == nil
}

// WriteAll 落盘整本缩略图（pngs 的 key 是 1 基页码）并写入有效性 stamp。
func WriteAll(dir, indexPath string, pngs map[int][]byte) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	// 先清旧页（上一次可能比这次多几页，删页后残留会展示已不存在的内容）
	ents, _ := os.ReadDir(dir)
	for _, e := range ents {
		if strings.HasSuffix(e.Name(), ".png") {
			_ = os.Remove(filepath.Join(dir, e.Name()))
		}
	}
	for no, png := range pngs {
		if len(png) == 0 {
			continue
		}
		if err := os.WriteFile(filepath.Join(dir, strconv.Itoa(no)+".png"), png, 0o644); err != nil {
			return err
		}
	}
	stamp := stampContent(indexPath)
	if stamp == "" {
		return fmt.Errorf("index.html 不存在，无法写缩略图 stamp")
	}
	return os.WriteFile(filepath.Join(dir, ".stamp"), []byte(stamp), 0o644)
}

// Invalidate 删除整份缩略图缓存（deck 页面任何写路径之后调用；不存在时静默）。
func Invalidate(decksDir, deckID string) {
	_ = os.RemoveAll(Dir(decksDir, deckID))
}

// Count 已缓存的缩略图页数（无缓存返回 0）。
func Count(dir string) int {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range ents {
		if strings.HasSuffix(e.Name(), ".png") {
			n++
		}
	}
	return n
}

// Pages 已缓存的页码（升序）。
func Pages(dir string) []int {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []int
	for _, e := range ents {
		name := strings.TrimSuffix(e.Name(), ".png")
		if n, err := strconv.Atoi(name); err == nil && n > 0 {
			out = append(out, n)
		}
	}
	sort.Ints(out)
	return out
}
