package usermodel

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"testing"

	"html-ppt/backend/internal/cryptox"
	"html-ppt/backend/internal/store"
)

// OpenMemory 是进程级共享库（cache=shared），用户 id 各测试独占段位，
// 且必须按传入 uid 建用户行（SetActive/ActiveTarget 都挂在 users 表上）。
func newTestService(t *testing.T, uid uint) (*Service, context.Context) {
	t.Helper()
	st, err := store.OpenMemory(t.Context())
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		t.Fatalf("rand: %v", err)
	}
	box, err := cryptox.NewBox(base64.StdEncoding.EncodeToString(raw))
	if err != nil {
		t.Fatalf("NewBox: %v", err)
	}
	u := store.User{Email: "um-" + base64.RawStdEncoding.EncodeToString(raw[:4]) + "@example.com"}
	u.ID = uid
	if err := st.DB.Create(&u).Error; err != nil {
		t.Fatalf("造用户: %v", err)
	}
	return New(st.DB, box), context.Background()
}

func TestCreateDefaultsAndSealsKey(t *testing.T) {
	uid := uint(95_001)
	s, ctx := newTestService(t, uid)

	row, err := s.Create(ctx, uid, "", "gpt-4o", "https://api.openai.com/v1/", "sk-test-123")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	// 名称默认 = 模型 ID；baseURL 去尾斜杠；key 已加密（不可逆读）
	if row.Name != "gpt-4o" || row.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("默认值不对: %+v", row)
	}
	if row.APIKeyEnc == "" || row.APIKeyEnc == "sk-test-123" {
		t.Errorf("key 必须加密落库: %q", row.APIKeyEnc)
	}

	// ActiveTarget 解出明文 key 与模型
	if err := s.SetActive(ctx, uid, row.ID); err != nil {
		t.Fatalf("SetActive: %v", err)
	}
	res, err := s.ActiveTarget(ctx, uid)
	if err != nil || res.Model == nil {
		t.Fatalf("ActiveTarget: %v %+v", err, res)
	}
	if res.Model.ModelID != "gpt-4o" || res.APIKey != "sk-test-123" {
		t.Errorf("解析目标不对: %+v key=%q", res.Model, res.APIKey)
	}
}

func TestActiveDefaultsToPlatformWhenNone(t *testing.T) {
	s, ctx := newTestService(t, 95_002)
	res, err := s.ActiveTarget(ctx, 95_002)
	if err != nil {
		t.Fatalf("ActiveTarget: %v", err)
	}
	if res.Model != nil {
		t.Errorf("没设过自选模型应回平台（nil），实际 %+v", res.Model)
	}
}

func TestUpdateKeepsKeyWhenEmpty(t *testing.T) {
	uid := uint(95_003)
	s, ctx := newTestService(t, uid)
	row, err := s.Create(ctx, uid, "我的模型", "m1", "https://a.example.com/v1", "sk-old")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := s.Update(ctx, uid, row.ID, "改名", "m2", "https://b.example.com/v1", ""); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if err := s.SetActive(ctx, uid, row.ID); err != nil {
		t.Fatalf("SetActive: %v", err)
	}
	res, err := s.ActiveTarget(ctx, uid)
	if err != nil || res.Model == nil {
		t.Fatalf("ActiveTarget: %v", err)
	}
	// 名字/模型/入口都更新了，key 保留旧的
	if res.Model.Name != "改名" || res.Model.ModelID != "m2" ||
		res.Model.BaseURL != "https://b.example.com/v1" || res.APIKey != "sk-old" {
		t.Errorf("Update 结果不对: %+v key=%q", res.Model, res.APIKey)
	}
}

func TestDeleteActiveResetsPointer(t *testing.T) {
	uid := uint(95_004)
	s, ctx := newTestService(t, uid)
	row, _ := s.Create(ctx, uid, "", "m1", "https://a.example.com/v1", "sk-1")
	if err := s.SetActive(ctx, uid, row.ID); err != nil {
		t.Fatalf("SetActive: %v", err)
	}
	if err := s.Delete(ctx, uid, row.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	// 行没了；指针归零（回平台模型）
	if _, err := s.owned(ctx, uid, row.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("删除后应 NotFound，实际 %v", err)
	}
	res, err := s.ActiveTarget(ctx, uid)
	if err != nil || res.Model != nil {
		t.Errorf("删除当前使用的模型后应回平台，实际 %+v err=%v", res, err)
	}
}

func TestOwnershipIsolation(t *testing.T) {
	uidA, uidB := uint(95_005), uint(95_006)
	s, ctx := newTestService(t, uidA)
	newTestService(t, uidB) // 只建 B 的用户行
	row, _ := s.Create(ctx, uidA, "", "m1", "https://a.example.com/v1", "")
	if err := s.SetActive(ctx, uidB, row.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("别人的模型不能设为当前: %v", err)
	}
	if err := s.Delete(ctx, uidB, row.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("别人的模型不能删: %v", err)
	}
}

func TestValidationAndLimit(t *testing.T) {
	uid := uint(95_007)
	s, ctx := newTestService(t, uid)

	if _, err := s.Create(ctx, uid, "", "", "https://a.com/v1", ""); !errors.Is(err, ErrModelIDRequired) {
		t.Errorf("空模型 ID 应拒绝: %v", err)
	}
	if _, err := s.Create(ctx, uid, "", "m1", "ftp://a.com", ""); !errors.Is(err, ErrBaseURLInvalid) {
		t.Errorf("非 http 入口应拒绝: %v", err)
	}
	// 上限 10 个
	for i := 0; i < MaxModelsPerUser; i++ {
		if _, err := s.Create(ctx, uid, "", "m", "https://a.com/v1", ""); err != nil {
			t.Fatalf("第 %d 个 Create: %v", i+1, err)
		}
	}
	if _, err := s.Create(ctx, uid, "", "m11", "https://a.com/v1", ""); !errors.Is(err, ErrLimitReached) {
		t.Errorf("第 11 个应拒绝: %v", err)
	}
}
