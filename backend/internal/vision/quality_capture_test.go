package vision

// 质量验收的截图工具：QUALITY_URL 指向（带鉴权的）deck 页面时，
// 全页截图落盘供人工逐页审查。本地：QUALITY_URL=... go test -run TestQualityCapture
import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestQualityCapture(t *testing.T) {
	if testing.Short() || os.Getenv("QUALITY_URL") == "" {
		t.Skip("设 QUALITY_URL 以运行质量截图")
	}
	d, err := CaptureV2(t.Context(), OptionsV2{URL: os.Getenv("QUALITY_URL")})
	if err != nil {
		t.Fatalf("捕获失败: %v", err)
	}
	outDir := os.Getenv("QUALITY_OUT")
	if outDir == "" {
		outDir = filepath.Join("..", "..", "..", "tmp", "quality")
	}
	_ = os.MkdirAll(outDir, 0o755)
	for _, s := range d.Slides {
		if len(s.PNG) == 0 {
			continue
		}
		p := filepath.Join(outDir, fmt.Sprintf("page-%02d.png", s.Index+1))
		if err := os.WriteFile(p, s.PNG, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("量测 %d 页:\n%s", len(d.Slides), DigestV2(d))
}
