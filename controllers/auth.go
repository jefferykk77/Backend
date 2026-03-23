package controllers

import (
	"ExchangeApp/global"
	"ExchangeApp/models"
	"ExchangeApp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	//models.User中不能用json标签控制Automigrate时创建的列名
	var user models.User

	//ctx.ShouldBindJSON将前端post的数据反序列化为结构体，赋值给user对象
	//赋值时先找json，没有就找字段，所以最好添加json，字段名匹配是大小写敏感的
	if err := ctx.ShouldBindJSON(&user); err != nil {
		//ctx.JSON创建了一个响应，包含状态码、响应体，自动添加相应头Content-Type: application/json。
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	hashedpwd, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	user.Password = hashedpwd

	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.AutoMigrate(&user); err != nil { //这个可以写在InitDB中
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})

}
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := global.Db.Where("Username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong cedentials"})
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong cedentials"})
		return
	}

	token, err := utils.GenerateJWT(input.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
