package controllers

import (
	"ExchangeApp/global"
	"ExchangeApp/models"
	"encoding/json"
	"errors"
	"net/http"
	"time"

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
	//exp过期或者要更新文章，可以采用清除缓存的方法
	if err := global.RedisDB.Del(cacheKey).Err(); err != nil { //把
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, article)
}

var cacheKey = "article"

func GetArticle(ctx *gin.Context) {
	//从缓存中寻找
	cacheData, err := global.RedisDB.Get(cacheKey).Result()
	//redis中不存在该键，从数据库中寻找
	if err == redis.Nil {
		var articles []models.Article
		//find找所有的，故用结构体切片
		//find 包含一步unmarshal吧
		if err := global.Db.Find(&articles).Error; err != nil { //db和redis db里的数据是字节流
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //.Error中的结构体赋值给err，err类型为error 接口类型，接口实例自然可调用接口方法
				//Error内有结构体，他实现了Error()方法，赋值err后，err中tab 存储了方法对应的真实底层函数地址，可直接调用，详情见GEM
			}
			return
		}
		cacheJson, err := json.Marshal(articles)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//get声明kv对和exp
		if err := global.RedisDB.Set(cacheKey, cacheJson, 24*time.Hour).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//set之后，用户无需刷新即可获取缓存中的数据故用ctx
		ctx.Data(http.StatusOK, "application/json", []byte(cacheJson))
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		var articles []models.Article
		//unmarshal若要转换为结构体类型则
		if err := json.Unmarshal([]byte(cacheData), &articles); err != nil { //字节流转换为第二个参数的数据类型，unmarshal第二个变量通常为指针，方便赋值
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, articles)
		//
		//Marshal（序列化）把 结构体/对象 变成 二进制/字符串
		//Unmarshal把 二进制/字符串 变回 结构体/对象（拆箱）
	}
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
	if err := global.RedisDB.Incr(ArticleLikesKey).Err(); err != nil { //需要获取
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
		ctx.JSON(http.StatusOK, gin.H{"likes": "0"})
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"likes": likes})
	}
}
