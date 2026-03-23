package router

import (
	"ExchangeApp/config"
	"ExchangeApp/controllers"
	"ExchangeApp/middlewears"

	"github.com/gin-gonic/gin"
)

// gin.Engine 是 Gin 框架的核心结构体，负责路由分发、中间件管理和服务启动。
func SetRouter() *gin.Engine {
	config.InitDB()
	config.InitConfig()
	r := gin.Default()

	//Group 用于路由分组，通常代表 URL 路径的共同前缀。
	//URL (Uniform Resource Locator): 统一资源定位符。俗称“网址”，用于在网络上定位资源的位置
	//HTTP (HyperText Transfer Protocol): 超文本传输协议。用于在客户端（浏览器/App）和服务器之间传输数据的规则。

	/*
		二者的关系：
			URL 是地址，HTTP 是搬运方式。
			URL 通常以 http:// 或 https:// 开头，这表示“通过 HTTP 协议去访问这个地址上的资源”。
	*/
	auth := r.Group("/api/auth")
	{
		//POST 是一种 HTTP 请求方法,通常用于向服务器提交数据（创建资源或提交表单）。
		//POST (Request): 前端->后端。前端将数据（如用户名、密码）放入 Request Body 中发送。
		//gin.H (Response): 后端->前端。后端处理完逻辑，将数据封装在 gin.H 中，通过 Response Body 返回。
		//响应体是 由 gin.H 转换而来的 JSON 数据。
		//Postman 里的表现：点击 "Send" 后，Postman 底部控制台显示的那个 {"msg": "Login Success"} 就是响应体。
		auth.POST("/login", controllers.Login)
		//Status: 设置 HTTP 状态码。例如 http.StatusOK (200) 表示请求成功。
		//JSON: 将后方的 gin.H（Map 别名）序列化为 JSON 格式并写入响应体。设置 Content-Type 为 application/json。
		/*
			Abort: 中断当前请求后续的所有处理流程。
				在 Gin 中，请求经过一系列中间件（Handler Chain）。
				调用 Abort 会确保后续的中间件或处理器不再执行，直接将当前结果返回给客户端

		*/

		auth.POST("/register", controllers.Register)

	}
	/*
		当你使用 Postman 访问 /api/auth/login 时，执行逻辑如下：

			Postman 发起请求：你点击 Send，发送一个 POST 请求到服务器。

			Gin 路由匹配：Gin 识别到路径是 /api/auth/login 且方法是 POST，进入你定义的匿名函数。

			后端处理 AbortWithStatusJSON：

			Status: 后端准备好 200 OK 状态码。

			JSON: 后端将 gin.H{"msg": "Login Success"} 转换（序列化）为字符串 {"msg": "Login Success"}。

			Abort: 告诉 Gin 引擎：“这个请求处理完了，不要再执行后面的逻辑了”。

			数据传输：后端将状态码 200 和 JSON 字符串通过 HTTP 协议发送回 Postman。

			Postman 接收展示：Postman 收到数据，显示状态码 200，并在 Body 区域渲染出 msg 内容。
	*/

	api := r.Group("/api")
	api.GET("/exchangeRates", controllers.GetExchangeRates)
	api.Use(middlewears.AuthMiddleWear())
	{
		api.POST("/exchangeRates", controllers.CreateExchangeRate)
	}
	return r
}
