package routers

import (
	"FoodBackend/middleware"
	api "FoodBackend/pkg/api/controllers"
	"FoodBackend/pkg/setting"
	"FoodBackend/routers/third_party/mini_app"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func InitRouter() *gin.Engine {
	fmt.Println("Hi")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	//r.Use(middleware.Secure)
	gin.SetMode(setting.RunMode)

	// swagger api docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// auth login
	r.POST("/api/auth", mini_app.GetAuth)

	// basic web info
	r.GET("/api/info", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": map[string]interface{}{"name": "嘻嘻嘻，小张私厨房app"},
		})
	})

	bookController := &api.BookController{}
	CommentController := &api.CommentController{}
	apiv1Guest := r.Group("/api")
	{
		apiv1Guest.GET("/books", bookController.List)
		apiv1Guest.GET("/books/:id", bookController.Get)

		apiv1Guest.GET("/comments", CommentController.List)
	}

	apiv1 := r.Group("/api", middleware.JWT())
	//apiv1 := r.Group("/api")
	{
		apiv1.POST("/books", bookController.Create)
		apiv1.PATCH("/books/status/:id", bookController.ChangeStatus)
		apiv1.PUT("/books/:id", bookController.Update)
		apiv1.DELETE("/books/:id", bookController.Delete)

		TagController := &api.TagController{}
		apiv1.GET("/tags", TagController.List)

		apiv1.POST("/comments", CommentController.Create)
		apiv1.DELETE("/comments/:id", CommentController.Delete)

		UserController := &api.UserController{}
		apiv1.GET("/user/info", UserController.Get)
	}

	return r
}
