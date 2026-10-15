package handler

import (
	"strings"
	"testing"
)

func TestExportDownloadName(t *testing.T) {
	cases := []struct {
		name     string
		title    string
		diskName string
		want     string
	}{
		{"中文标题原样", "人工智能入门", "deck.html", "人工智能入门.html"},
		{"非法字符剔除", `a/b\c:d*e?f"g<h>i|j`, "deck.pdf", "abcdefghij.pdf"},
		{"控制符剔除（含 tab）", "  题\x00\t目  ", "deck.pdf", "题目.pdf"},
		{"超长按 rune 截断到 60", strings.Repeat("题", 100), "deck.html", strings.Repeat("题", 60) + ".html"},
		{"全空白回退 deck", "  /?  ", "deck-png.zip", "deck.zip"},
		{"扩展名跟随产物", "周报", "deck-png.zip", "周报.zip"},
		{"首尾点空格修剪", "  .标题.  ", "deck.pdf", "标题.pdf"},
	}
	for _, tc := range cases {
		if got := exportDownloadName(tc.title, tc.diskName); got != tc.want {
			t.Errorf("%s: exportDownloadName(%q, %q) = %q, want %q", tc.name, tc.title, tc.diskName, got, tc.want)
		}
	}
}
