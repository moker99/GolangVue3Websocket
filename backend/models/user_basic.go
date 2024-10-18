package models

import (
	"fmt"
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name                string
	PassWord            string
	Phone               string `valid:"matches(^09[0-9]{8}$)"`
	Email               string `valid:"email"`
	Identity            string
	ClientIp            string
	ClientPort          string
	Salt                string
	LoginTime           *time.Time
	HeartbeatTime       *time.Time
	LoginOutTime        *time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout            bool
	DeviceInfo          string
	LatestMessage       string `gorm:"-"` // 不映射到數據庫，用于存儲最新的消息
	LatestMessageTime   string `gorm:"-"` // 不映射到數據庫，用于存儲最新的消息的時間
	CountUnreadMessages int    `gorm:"-"` // 不映射到數據庫，用于紀錄未讀訊息數量
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func (user *UserBasic) BeforeSave(tx *gorm.DB) (err error) {
	user.LoginTime = nil
	user.HeartbeatTime = nil
	user.LoginOutTime = nil
	return
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and pass_word = ?", name, password).First(&user)

	// token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

// ///////////////////尚未使用////////////////////
func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone = ?", phone).First(&user)
}

func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email = ?", email).First(&user)
}

func FindByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id = ?", id).First(&user)
	return user
}

////////////////////////////////////////////////

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}
