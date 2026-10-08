package usermodel

// 用户自选模型：OpenAI 兼容入口（base_url + model_id + api_key）的增删改查、
// 「当前使用」指针，以及给 agent 解析调用目标。
//
// 选中后全部 LLM 链路改走用户的模型和 key——费用记用户自己的服务商，与平台
// 额度无关（平台本就没有按 key 的额度闸门，用量账本照记，只是 model 列变成
// 用户的模型 ID，用量页可按它筛选）。
//
// key 用与 BYOK 同一套 AES-GCM（cryptox.Box）加密落库，永不出后端：列表只回
// has_key 布尔，编辑不回填，测试连通由后端拿解密后的 key 代跑。

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
	"html-ppt/backend/internal/cryptox"
	"html-ppt/backend/internal/store"
)

const (
	MaxModelsPerUser = 10
	maxKeyLen        = 512
	maxNameLen       = 128
	maxModelIDLen    = 128
	maxBaseURLLen    = 255
)

var (
	ErrStorageUnavailable = errors.New("存储不可用")
	ErrCryptoUnavailable  = errors.New("加密组件不可用")
	ErrNotFound           = errors.New("模型不存在")
	ErrLimitReached       = errors.New("最多添加 10 个模型")
	ErrModelIDRequired    = errors.New("模型 ID 不能为空")
	ErrBaseURLInvalid     = errors.New("基础地址必须以 http:// 或 https:// 开头")
	ErrKeyInvalid         = errors.New("API Key 过长")
	ErrNameInvalid        = errors.New("模型名称过长")
	ErrModelIDInvalid     = errors.New("模型 ID 过长")
)

// Service 模型管理。db/box 任一为 nil 时所有方法返回对应不可用错误
// （数据库降级 / 未配 CRYPTO_AES_KEY 的部署，与 auth.Service 同款语义）。
type Service struct {
	db  *gorm.DB
	box *cryptox.Box
}

func New(db *gorm.DB, box *cryptox.Box) *Service {
	return &Service{db: db, box: box}
}

// Available 报告模型管理是否可用（handler 据此返回 503）。
func (s *Service) Available() bool { return s != nil && s.db != nil && s.box != nil }

// Resolved 解密后的调用目标。Model 为 nil 表示平台模型（调用方回退现状）。
type Resolved struct {
	Model  *store.UserModel
	APIKey string
}

// List 用户自己的模型清单（按添加时间倒序，新者在先）。
func (s *Service) List(ctx context.Context, uid uint) ([]store.UserModel, error) {
	if !s.Available() {
		return nil, s.unavailable()
	}
	var rows []store.UserModel
	err := s.db.WithContext(ctx).Where("user_id = ?", uid).
		Order("created_at DESC, id DESC").Find(&rows).Error
	return rows, err
}

// Create 新增。name 为空时取 modelID（用户拍板的默认规则）；apiKey 可空
// （本地网关常见）。每人最多 MaxModelsPerUser 个。
func (s *Service) Create(ctx context.Context, uid uint, name, modelID, baseURL, apiKey string) (*store.UserModel, error) {
	if !s.Available() {
		return nil, s.unavailable()
	}
	name, modelID, baseURL, err := normalize(name, modelID, baseURL)
	if err != nil {
		return nil, err
	}
	if apiKey = strings.TrimSpace(apiKey); len(apiKey) > maxKeyLen {
		return nil, ErrKeyInvalid
	}
	var cnt int64
	if err := s.db.WithContext(ctx).Model(&store.UserModel{}).Where("user_id = ?", uid).Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt >= MaxModelsPerUser {
		return nil, ErrLimitReached
	}
	enc := ""
	if apiKey != "" {
		if enc, err = s.box.Seal(apiKey); err != nil {
			return nil, err
		}
	}
	row := store.UserModel{
		UserID: uid, Name: name, ModelID: modelID, BaseURL: baseURL, APIKeyEnc: enc,
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Update 改。apiKey 为空 = 保留原 key（编辑时不回填、不改传空）。
func (s *Service) Update(ctx context.Context, uid, id uint, name, modelID, baseURL, apiKey string) error {
	if !s.Available() {
		return s.unavailable()
	}
	name, modelID, baseURL, err := normalize(name, modelID, baseURL)
	if err != nil {
		return err
	}
	if apiKey = strings.TrimSpace(apiKey); len(apiKey) > maxKeyLen {
		return ErrKeyInvalid
	}
	row, err := s.owned(ctx, uid, id)
	if err != nil {
		return err
	}
	row.Name, row.ModelID, row.BaseURL = name, modelID, baseURL
	if apiKey != "" {
		enc, err := s.box.Seal(apiKey)
		if err != nil {
			return err
		}
		row.APIKeyEnc = enc
	}
	return s.db.WithContext(ctx).Model(row).Updates(map[string]any{
		"name": row.Name, "model_id": row.ModelID, "base_url": row.BaseURL, "api_key_enc": row.APIKeyEnc,
	}).Error
}

// Delete 删。删的是「当前使用」的模型则指针归零（切回平台模型）。
func (s *Service) Delete(ctx context.Context, uid, id uint) error {
	if !s.Available() {
		return s.unavailable()
	}
	if _, err := s.owned(ctx, uid, id); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND user_id = ?", id, uid).Delete(&store.UserModel{}).Error; err != nil {
			return err
		}
		return tx.Model(&store.User{}).Where("id = ? AND active_model_id = ?", uid, id).
			Update("active_model_id", 0).Error
	})
}

// SetActive 设「当前使用」。id=0 = 切回平台模型；id 必须属于该用户。
func (s *Service) SetActive(ctx context.Context, uid, id uint) error {
	if !s.Available() {
		return s.unavailable()
	}
	if id != 0 {
		if _, err := s.owned(ctx, uid, id); err != nil {
			return err
		}
	}
	return s.db.WithContext(ctx).Model(&store.User{}).
		Where("id = ?", uid).Update("active_model_id", id).Error
}

// ActiveTarget 解析当前该用哪个模型：active_model_id 指向的自选模型（解密 key），
// 没有指向 / 行不存在 / 解密失败 → nil（调用方回退平台模型，fail-open 同 BYOK）。
func (s *Service) ActiveTarget(ctx context.Context, uid uint) (*Resolved, error) {
	if !s.Available() {
		return nil, s.unavailable()
	}
	var u store.User
	if err := s.db.WithContext(ctx).Model(&store.User{}).Where("id = ?", uid).
		Select("id", "active_model_id").First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &Resolved{}, nil
		}
		return nil, err
	}
	if u.ActiveModelID == 0 {
		return &Resolved{}, nil
	}
	var row store.UserModel
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", u.ActiveModelID, uid).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &Resolved{}, nil // 指针悬空（理论上有 Delete 兜底）：按平台模型走
		}
		return nil, err
	}
	key := ""
	if row.APIKeyEnc != "" {
		var err error
		if key, err = s.box.Open(row.APIKeyEnc); err != nil {
			// 解密失败（换过 CRYPTO_AES_KEY 之类）：宁可回退平台模型也别把对话打断
			return &Resolved{}, nil
		}
	}
	return &Resolved{Model: &row, APIKey: key}, nil
}

// Test 连通性：对给定入口发一次最小的非流式补全（15s 超时）。
// 不设 MaxTokens/ReasoningEffort——与全仓库"上限交给默认值"的教训一致，
// 推理模型把小上限吃满只会 finish=length，调用本身仍算通。
func (s *Service) Test(ctx context.Context, baseURL, modelID, apiKey string) error {
	if !s.Available() {
		return s.unavailable()
	}
	_, _, _, err := normalize("", modelID, baseURL)
	if err != nil {
		return err
	}
	client := NewClient(baseURL, apiKey, 0)
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	completion, err := client.Chat.Completions.New(ctx, newTestRequest(modelID))
	if err != nil {
		return err
	}
	if len(completion.Choices) == 0 {
		return errors.New("响应没有 choices")
	}
	return nil
}

// owned 归属校验（不存在与别人的返回同一个错误，不泄露存在性）。
func (s *Service) owned(ctx context.Context, uid, id uint) (*store.UserModel, error) {
	var row store.UserModel
	err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, uid).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (s *Service) unavailable() error {
	if s == nil || s.db == nil {
		return ErrStorageUnavailable
	}
	return ErrCryptoUnavailable
}

func normalize(name, modelID, baseURL string) (string, string, string, error) {
	name = strings.TrimSpace(name)
	modelID = strings.TrimSpace(modelID)
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if modelID == "" {
		return "", "", "", ErrModelIDRequired
	}
	if len(modelID) > maxModelIDLen {
		return "", "", "", ErrModelIDInvalid
	}
	if len(name) > maxNameLen {
		return "", "", "", ErrNameInvalid
	}
	if name == "" {
		name = modelID // 用户拍板：名称默认用模型 ID
	}
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		return "", "", "", ErrBaseURLInvalid
	}
	if len(baseURL) > maxBaseURLLen {
		return "", "", "", ErrBaseURLInvalid
	}
	return name, modelID, baseURL, nil
}

// HTTPStatus 服务错误 → 响应码（handler 用）。
func HTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrStorageUnavailable), errors.Is(err, ErrCryptoUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrLimitReached), errors.Is(err, ErrModelIDRequired),
		errors.Is(err, ErrBaseURLInvalid), errors.Is(err, ErrKeyInvalid),
		errors.Is(err, ErrNameInvalid), errors.Is(err, ErrModelIDInvalid):
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

// ConfigFor 取某行解密后的调用配置（连通性测试用）；key 未设时返回空串。
func (s *Service) ConfigFor(ctx context.Context, uid, id uint) (baseURL, modelID, apiKey string, err error) {
	if !s.Available() {
		return "", "", "", s.unavailable()
	}
	row, err := s.owned(ctx, uid, id)
	if err != nil {
		return "", "", "", err
	}
	if row.APIKeyEnc != "" {
		if apiKey, err = s.box.Open(row.APIKeyEnc); err != nil {
			return "", "", "", err
		}
	}
	return row.BaseURL, row.ModelID, apiKey, nil
}

// ActiveID 当前选中的自选模型 id（0 = 平台模型）。列表页展示「当前使用」用。
func (s *Service) ActiveID(ctx context.Context, uid uint) (uint, error) {
	if !s.Available() {
		return 0, s.unavailable()
	}
	var u store.User
	if err := s.db.WithContext(ctx).Model(&store.User{}).Where("id = ?", uid).
		Select("id", "active_model_id").First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return u.ActiveModelID, nil
}
