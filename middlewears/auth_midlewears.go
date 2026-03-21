package middlewears

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWear() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		token:=ctx.GetHeader("Authorization")
		if token==""{
			ctx.JSON(http.StatusUnauthorized,gin.H{"error":"Missing Authorization Header"})
			ctx.Abort()
			return 
		}
	}
}