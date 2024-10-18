package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList
// @Summary 所有用戶
// @Tags 用戶模塊
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(200, gin.H{
		"code":    0, //  0成功  -1失敗
		"message": "查詢成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用戶
// @Tags 用戶模塊
// @param name query string false "用戶名"
// @param password query string false "密碼"
// @param repassword query string false "確認密碼"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")

	salt := fmt.Sprintf("%06d", rand.Int31())

	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功  -1失敗
			"message": "用戶名或密碼不能為空！！",
			"data":    user,
		})
		return
	}
	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功  -1失敗
			"message": "用戶名重複註冊",
			"data":    user,
		})
		return
	}
	if password != repassword {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功  -1失敗
			"message": "兩次密碼不一致！！",
			"data":    user,
		})
		return
	}

	// token加密
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt

	models.CreateUser(user)
	data = models.FindUserByNameAndPwd(user.Name, user.PassWord)

	c.JSON(200, gin.H{
		"code":    0, //  0成功  -1失敗
		"message": "新增用戶成功！！",
		"data":    data,
	})
}

// FindUserByNameAndPwd
// @Summary 用戶登入
// @Tags 用戶模塊
// @param name query string false "用戶名"
// @param password query string false "密碼"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功  -1失敗
			"message": "該用戶不存在",
			"data":    data,
		})
		return
	}

	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功  -1失敗
			"message": "密碼不正確",
			"data":    data,
		})
		return
	}

	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)

	c.JSON(200, gin.H{
		"code":    0, //  0成功  -1失敗
		"message": "登入成功",
		"data":    data,
	})
}

// DeleteUser
// @Summary 刪除用戶
// @Tags 用戶模塊
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, //  0成功  -1失敗
		"message": "刪除用戶成功！！",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用戶
// @Tags 用戶模塊
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, //  0成功  -1失敗
			"message": "修改參數不匹配！！",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0, //  0成功  -1失敗
			"message": "修改用戶成功！！",
			"data":    user,
		})
	}

}

// 防止跨域站點偽造請求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("發送消息： ", msg)
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func LoadFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.LoadFriends(uint(id))

	utils.RespOKList(c.Writer, users, len(users))
}

func AddFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	friendName := c.Request.FormValue("friendName")
	code, msg := models.AddFriend(uint(id), friendName)

	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	desc := c.Request.FormValue("desc")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Desc = desc
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))

	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func GetChatHistory(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.Request.FormValue("userIdA"))
	userIdB, _ := strconv.Atoi(c.Request.FormValue("userIdB"))
	start, _ := strconv.Atoi(c.Request.FormValue("start"))
	end, _ := strconv.Atoi(c.Request.FormValue("end"))
	res, err := models.GetChatHistory(int(userIdA), int(userIdB), int64(start), int64(end))
	if err != nil {
		fmt.Println(err)
	}
	utils.RespOKList(c.Writer, res, len(res))
}
