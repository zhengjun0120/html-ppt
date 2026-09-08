package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html-ppt/backend/internal/store"
	"strings"

	"github.com/openai/openai-go/v3"
	"gorm.io/gorm"
)

// 控制信号，循环在ask_user处暂停
var ErrPaused = errors.New("agent paused: waiting for user answer")

// 续轮读回历史，首轮新建会话
func (as *AgentService) loadOrCreateSession(ctx context.Context,userID,sessionID uint,userContent,deckID string) (*store.ChatSession,[]openai.ChatCompletionMessageParamUnion,error){
	if as.st == nil{
		return nil,nil,errors.New("数据库不可用，会话功能关闭")
	}

	// 有sessionID 就是可以恢复会话
	if sessionID > 0{
		var sess store.ChatSession
		err := as.st.DB.Where("id = ? AND user_id = ?",sessionID,userID).First(&sess).Error
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil,fmt.Errorf("会话不存在")
		}
		if err !=nil{
			return nil,nil,fmt.Errorf("加载会话时 数据库失败 err:%w",err)
		}

		var messages []openai.ChatCompletionMessageParamUnion
		if sess.Messages != ""{
			if err := json.Unmarshal([]byte(sess.Messages),&messages);err !=nil{
				return nil,nil,fmt.Errorf("会话消息损坏 err:%w",err)
			}
		}
		return &sess,messages,nil
	}

	// 没有就是首轮 首轮需要入库
	title := strings.TrimSpace(userContent)
	if r := []rune(title);len(r) >30{
		title = string(r[:30])
	}
	sess := &store.ChatSession{UserID: userID,DeckID: deckID,Title:title}
	if err := as.st.DB.Create(sess).Error;err!=nil{
		return nil,nil,fmt.Errorf("创建会话时 数据库失败 err:%w",err)
	}
	return sess,nil,nil
}

// 仅读回 （恢复入口用，不新建会话）
func (as *AgentService) loadSession(ctx context.Context,userID,sessionID uint)(*store.ChatSession,[]openai.ChatCompletionMessageParamUnion,error){
	if sessionID == 0{
		return nil,nil,errors.New("缺少 session_id")
	}
	sess,messages,err := as.loadOrCreateSession(ctx,userID,sessionID,"","")
	return sess,messages,err
}

// 把完整消息数组写回会话，每个检查点整体重写
func (as *AgentService) persistSession(sess *store.ChatSession,messages []openai.ChatCompletionMessageParamUnion) error{
	data,err := json.Marshal(messages)
	if err != nil{
		return fmt.Errorf("序列化会话消息失败 err:%w",err)
	}
	sess.Messages = string(data)
	err = as.st.DB.Model(sess).Updates(map[string]any{
		"messages":sess.Messages,
		"pending_ask":sess.PendingAsk,
	}).Error
	if err !=nil{
		return fmt.Errorf("更新会话时 数据库时失败 err:%w",err)
	}
	return nil
}

// 预检验 查询会话是否存在
func (as *AgentService) EnsureSessionOwner(userID,sessionID uint) error{
	if as.st == nil{
		return errors.New("数据库不可用")
	}
	var cnt int64
	if err := as.st.DB.Model(&store.ChatSession{}).Where("id = ? AND user_id = ?",sessionID,userID).Count(&cnt).Error;err !=nil{
		return fmt.Errorf("查询会话: %w",err)
	}
	if cnt == 0{
		return fmt.Errorf("会话不存在")
	}
	return nil
}