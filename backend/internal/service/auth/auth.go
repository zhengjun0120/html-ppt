// Package auth 负责注册/登录/验证码/用户自带 API Key 的领域逻辑。
// 身份只有邮箱一个维度；所有秘密（密码、API Key、验证码）都不以明文落库/落日志。
package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/cryptox"
	"html-ppt/backend/internal/store"
)

// 哨兵错误：handler 层用 errors.Is 映射 HTTP 状态码。
// 文案刻意不带"用户是否存在"之类信息，防账号枚举。
var (
	ErrEmailInvalid       = errors.New("邮箱格式不正确")
	ErrPasswordWeak       = errors.New("密码需 8~72 位")
	ErrAPIKeyInvalid      = errors.New("API Key 格式不正确")
	ErrCodeInvalid        = errors.New("验证码错误或已过期")
	ErrCodeCooldown       = errors.New("发送太频繁，请 1 分钟后再试")
	ErrEmailTaken         = errors.New("该邮箱已注册")
	ErrBadCredentials     = errors.New("邮箱或密码不正确")
	ErrTokenRefreshQuota  = errors.New("今天已经刷新过登录凭据")
	ErrCryptoUnavailable  = errors.New("加密模块未配置（crypto.aes_key），无法保存 API Key")
	ErrStorageUnavailable = errors.New("数据库不可用")
)

const (
	codeKey      = "auth:code:"    // + email → 验证码，TTL 10 分钟
	attemptsKey  = "auth:code:try:" // + email → 验证码尝试次数，TTL 同验证码
	cooldownK    = "auth:code:cd:" // + email → 发送冷却标记，TTL 60 秒
	codeTTL      = 10 * time.Minute
	cooldown     = time.Minute
	minPWLen     = 8
	maxPWLen     = 72  // bcrypt 算法硬上限：超过部分被静默忽略，必须显式拒绝
	maxEmailLen  = 254 // RFC 5321 上限，也与 users.email 列宽匹配
	maxAPIKeyLen = 512
	maxCodeAttempts = 5 // 一个验证码最多试 5 次，防 6 位码被 10 分钟内暴力穷举
)

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// validEmail = 长度上限 + 格式。长度必须显式卡：超长邮箱会撞列宽报 500。
func validEmail(e string) bool {
	return len(e) <= maxEmailLen && emailRe.MatchString(e)
}

type Service struct {
	db     *gorm.DB
	rdb    *redis.Client
	box    *cryptox.Box // 可为 nil（未配置 aes_key）：登录注册可用，API Key 功能报 503
	ttl    time.Duration
	secret []byte
	smtp   config.SMTP
}

func New(db *gorm.DB, rdb *redis.Client, box *cryptox.Box, secret string, ttl time.Duration, smtpCfg config.SMTP) *Service {
	return &Service{db: db, rdb: rdb, box: box, ttl: ttl, secret: []byte(secret), smtp: smtpCfg}
}

// RequestCode 给邮箱发 6 位注册验证码。两道闸：格式 → 冷却。
func (s *Service) RequestCode(ctx context.Context, email string) error {
	email = normalizeEmail(email)
	if !validEmail(email) {
		return ErrEmailInvalid
	}
	if s.rdb == nil {
		return ErrStorageUnavailable
	}

	// SetNX 天然原子：拿不到 = 冷却期内
	ok, err := s.rdb.SetNX(ctx, cooldownK+email, 1, cooldown).Result()
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	if !ok {
		return ErrCodeCooldown
	}

	code, err := genCode()
	if err != nil {
		return err
	}
	if err := s.rdb.Set(ctx, codeKey+email, code, codeTTL).Err(); err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	if err := s.sendCodeMail(email, code); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	return nil
}

// Register 验证码通过后建号，返回新用户 id。
func (s *Service) Register(ctx context.Context, email, code, password string) (uint, error) {
	email = normalizeEmail(email)
	if !validEmail(email) {
		return 0, ErrEmailInvalid
	}
	if len(password) < minPWLen || len(password) > maxPWLen {
		return 0, ErrPasswordWeak
	}
	if s.rdb == nil || s.db == nil {
		return 0, ErrStorageUnavailable
	}

	// 尝试次数限流：每个验证码最多试 5 次，超限作废（防 6 位码被穷举）。
	// Incr 是原子的；首次自增顺手设置过期，与验证码同生命周期。
	tries, err := s.rdb.Incr(ctx, attemptsKey+email).Result()
	if err != nil {
		return 0, fmt.Errorf("redis: %w", err)
	}
	if tries == 1 {
		s.rdb.Expire(ctx, attemptsKey+email, codeTTL)
	}
	if tries > maxCodeAttempts {
		s.rdb.Del(ctx, codeKey+email) // 锁死：当前验证码作废，必须重新申请
		return 0, ErrCodeInvalid
	}

	saved, err := s.rdb.Get(ctx, codeKey+email).Result()
	if errors.Is(err, redis.Nil) {
		return 0, ErrCodeInvalid
	}
	if err != nil {
		return 0, fmt.Errorf("redis: %w", err)
	}
	if !cryptox.ConstantTimeEqual(saved, strings.TrimSpace(code)) {
		return 0, ErrCodeInvalid
	}
	// 一次性使用：无论注册成败都作废，防验证码被反复试
	if err := s.rdb.Del(ctx, codeKey+email, attemptsKey+email).Err(); err != nil {
		return 0, fmt.Errorf("redis: %w", err)
	}

	var cnt int64
	if err := s.db.Model(&store.User{}).Where("email = ?", email).Count(&cnt).Error; err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}
	if cnt > 0 {
		return 0, ErrEmailTaken
	}

	hash, err := cryptox.HashPassword(password)
	if err != nil {
		return 0, err
	}
	u := store.User{Email: email, PasswordHash: hash}
	if err := s.db.Create(&u).Error; err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}
	return u.ID, nil
}

// Login 校验密码并签发 JWT。用户不存在和密码错误返回同一个错误、
// 耗时也对齐（都做一次 bcrypt），不给攻击者侧信息。
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	email = normalizeEmail(email)
	if s.db == nil {
		return "", ErrStorageUnavailable
	}

	// 静默 SQL 日志：登录路径不该把参数打进日志
	var u store.User
	err := s.db.Session(&gorm.Session{Logger: logger.Discard}).
		Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 与真实路径等价的 bcrypt 开销，抹平"存在与否"的响应时延差
		bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return "", ErrBadCredentials
	}
	if err != nil {
		return "", fmt.Errorf("db: %w", err)
	}
	if !cryptox.CheckPassword(u.PasswordHash, password) {
		return "", ErrBadCredentials
	}
	return s.signToken(u.ID), nil
}

// Refresh 给当前登录用户换发一个新 token（有效期从现在起重新计满一个 TTL）。
//
// 限额是"每用户每天一次"，额度记在用户行上、与设备无关：先上线的设备把当天的
// 额度用掉，其余设备拿着旧 token 继续用——JWT 无状态，已签发的都有效到过期，
// 所以"没抢到今天额度"不构成任何登录障碍，只是拿不到新 token 而已。
//
// 闸门是条件 UPDATE 的 RowsAffected，天然原子：并发两次刷新只有一次能把
// token_refreshed_at 从"今天之前"改成"现在"，另一次改 0 行、拿到额度用完的错。
// "今天"以服务器本地日历日为准（见 dayStart）。
func (s *Service) Refresh(ctx context.Context, userID uint) (string, error) {
	if s.db == nil {
		return "", ErrStorageUnavailable
	}
	now := time.Now()
	res := s.db.WithContext(ctx).Model(&store.User{}).
		Where("id = ? AND (token_refreshed_at IS NULL OR token_refreshed_at < ?)", userID, dayStart(now)).
		Update("token_refreshed_at", now)
	if res.Error != nil {
		return "", fmt.Errorf("db: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return "", ErrTokenRefreshQuota
	}
	return s.signToken(userID), nil
}

// dayStart 当天零点（本地时区）。"每天一次"的"天"以服务器本地日历日为准。
//
// 刻意不用 time.Truncate(24h)：它按"距零值时间的绝对时长"切分，切点落在 UTC 零点，
// 在东八区就是早上 8 点——那样"每天"会变成"每天早上 8 点到次日早上 8 点"。
func dayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// SetAPIKey 保存/清除（传空串清除）用户自带的 LLM API Key。存密文。
func (s *Service) SetAPIKey(ctx context.Context, userID uint, apiKey string) error {
	if s.db == nil {
		return ErrStorageUnavailable
	}
	if s.box == nil {
		return ErrCryptoUnavailable
	}

	apiKey = strings.TrimSpace(apiKey)
	enc := ""
	if apiKey != "" {
		if len(apiKey) > maxAPIKeyLen {
			return ErrAPIKeyInvalid
		}
		var err error
		enc, err = s.box.Seal(apiKey)
		if err != nil {
			return fmt.Errorf("加密 api key: %w", err)
		}
	}
	return s.db.Model(&store.User{}).Where("id = ?", userID).
		Update("api_key_enc", enc).Error
}

// Me 返回用户摘要（邮箱 + 是否配置了自带 key）。
func (s *Service) Me(ctx context.Context, userID uint) (email string, hasKey bool, err error) {
	if s.db == nil {
		return "", false, ErrStorageUnavailable
	}
	var u store.User
	if err := s.db.Select("email", "api_key_enc").First(&u, userID).Error; err != nil {
		return "", false, fmt.Errorf("db: %w", err)
	}
	return u.Email, u.APIKeyEnc != "", nil
}

// Token 为指定用户签发 JWT（注册成功后免二次登录直接进入会话）。
func (s *Service) Token(userID uint) string { return s.signToken(userID) }

func (s *Service) signToken(userID uint) string {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		panic("签发 JWT 失败: " + err.Error()) // 密钥格式错误属启动期配置错误，不应静默
	}
	return tok
}

func normalizeEmail(e string) string { return strings.ToLower(strings.TrimSpace(e)) }

func genCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
