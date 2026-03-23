package controllers

import (
	"ExchangeApp/global"
	"ExchangeApp/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func CreateArticle(ctx *gin.Context) {
	var article models.Article
	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.AutoMigrate(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, article)
}

func GetArticle(ctx *gin.Context) {
	var articles []models.Article
	if err := global.Db.Find(&articles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //.Error中的结构体赋值给err，err类型为error 接口类型，接口实例自然可调用接口方法
			//Error内有结构体，他实现了Error()方法，赋值err后，err中tab 存储了方法对应的真实底层函数地址，可直接调用
		}
		return
	}
	ctx.JSON(http.StatusOK, articles)
}

func GetCtxByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var article models.Article
	if err := global.Db.Where("id = ?", id).First(&article).Error; err != nil {
		if b := errors.Is(err, gorm.ErrRecordNotFound); b {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, &article)
}

func LikeArticle(ctx *gin.Context) {
	id := ctx.Param("id")
	ArticleLikesKey := "article:" + id + ":likes"
	if err := global.RedisDB.Incr(ArticleLikesKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"mssage": "successful"})
}

func GetArticleLikes(ctx *gin.Context) {
	id := ctx.Param("id")
	ArticleLikeKey := "article:" + id + ":likes"
	likes, err := global.RedisDB.Get(ArticleLikeKey).Result() //Redis `GET key` command. It returns redis.Nil error when key does not exist.
	if err == redis.Nil {
		likes = "0"
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"likes": likes})
	}
}
