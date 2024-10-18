package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

func CreateCommunity(community Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群名稱不能為空"
	}
	if community.OwnerId == 0 {
		return -1, "請先登入"
	}
	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		return -1, "建群失敗"
	}
	return 0, "建群成功"
}

func LoadCommunity(ownerId uint) ([]*Community, string) {
	data := make([]*Community, 10)
	utils.DB.Where("owner_id = ? ", ownerId).Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}

	return data, "查詢成功"
}
