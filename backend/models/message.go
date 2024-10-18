package models

import (
	"context"
	"encoding/json"
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

// 消息
type Message struct {
	gorm.Model
	FromId     int64
	TargetId   int64
	Type       int
	Media      int
	Content    string
	Pic        string
	Url        string
	Desc       string
	Amount     int
	CreateTime string `gorm:"-"`
	IsRead     int    `gorm:"-"`
}

func (table *Message) TableName() string {
	return "message"
}

func GetChatHistory(userIdA int, userIdB int, start int64, end int64) ([]string, error) {
	ctx := context.Background()
	var key string
	if userIdA > userIdB {
		key = fmt.Sprintf("msg_%d_%d_1", userIdB, userIdA)
	} else {
		key = fmt.Sprintf("msg_%d_%d_1", userIdA, userIdB)
	}

	messages, err := utils.Red.ZRange(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func GetLatestMessage(userIdA int, userIdB int) (Message, error) {
	ctx := context.Background()
	var key string
	if userIdA > userIdB {
		key = fmt.Sprintf("msg_%d_%d_1", userIdB, userIdA)
	} else {
		key = fmt.Sprintf("msg_%d_%d_1", userIdA, userIdB)
	}

	result, err := utils.Red.ZRange(ctx, key, -1, -1).Result()
	if err != nil {
		return Message{}, err
	}

	if len(result) == 0 {
		return Message{}, fmt.Errorf("no message found")
	}

	var message Message
	err = json.Unmarshal([]byte(result[0]), &message)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}
