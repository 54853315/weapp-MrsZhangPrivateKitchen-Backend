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
)

func InitRouter() *gin.Engine {
	fmt.Println("Hi")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	//r.Use(middleware.Cors())
	//r.Use(middleware.Secure)
	gin.SetMode(setting.RunMode)

	// swagger api docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// auth login
	r.POST("/api/auth", mini_app.GetAuth)

	apiv1 := r.Group("/api", middleware.JWT())
	{
		bookController := &api.BookController{}
		apiv1.GET("/books", bookController.List)
		apiv1.GET("/books/:id", bookController.Get)
		apiv1.POST("/books", bookController.Create)
		apiv1.PATCH("/books/status/:id", bookController.ChangeStatus)
		apiv1.PUT("/books/:id", bookController.Update)
		apiv1.DELETE("/books/:id", bookController.Delete)

		TagController := &api.TagController{}
		apiv1.GET("/tags", TagController.List)

		CommentController := &api.CommentController{}
		apiv1.GET("/comments", CommentController.List)
		apiv1.POST("/comments", CommentController.Create)
		apiv1.DELETE("/comments/:id", CommentController.Delete)
	}

	return r
}