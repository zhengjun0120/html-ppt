package handler

// mapDeckErr 的状态码语义守卫：阶段守卫拒绝（重复确认/重复提交）必须是 409，
// 不是 400 兜底——前端靠状态码区分"重复点击"和"请求本身坏了"。

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/service/deck"
)

func TestMapDeckErrStageMismatchIsConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mapDeckErr(c, deck.StageMismatch{Current: "selecting_template", Require: "outline_review"})
	if w.Code != http.StatusConflict {
		t.Fatalf("StageMismatch 应映射 409，实际 %d", w.Code)
	}
}
