package service

import (
	"ginchat/models"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首頁
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
}

func ToChat(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
}

func Chat(c *gin.Context) {
	models.UpgradeWebSocket(c.Writer, c.Request)
}
