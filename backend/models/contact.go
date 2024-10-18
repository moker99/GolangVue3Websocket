package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func LoadFriends(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts)

	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}

	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)

	for i, friend := range users {
		friendId := int(friend.ID)
		latestMessage, err := GetLatestMessage(int(userId), friendId)
		if err != nil {
			fmt.Println("Error getting latest message:", err)
			continue
		}

		users[i].LatestMessage = latestMessage.Content
		users[i].LatestMessageTime = latestMessage.CreateTime
	}

	return users
}

func AddFriend(userId uint, friendName string) (int, string) {
	user := UserBasic{}
	if friendName != "" {
		user = FindUserByName(friendName)
		if user.Identity != "" {
			if user.ID == userId {
				return -1, "不能加自己"
			}
			contact0 := Contact{}
			utils.DB.Where("owner_id = ? and target_id = ? and type = 1", userId, user.ID).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "該用戶已添加過!"
			}

			tx := utils.DB.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			contact := Contact{
				OwnerId:  userId,
				TargetId: user.ID,
				Type:     1,
			}
			if err := tx.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "添加失敗!"
			}

			contact2 := Contact{
				OwnerId:  user.ID,
				TargetId: userId,
				Type:     1,
			}
			if err := tx.Create(&contact2).Error; err != nil {
				tx.Rollback()
				return -1, "添加失敗!"
			}
			tx.Commit()

			return 0, "添加好友成功！"
		}
		return -1, "該用戶不存在"
	}

	return -1, "該用戶不存在"
}
