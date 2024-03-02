package routes

import (
	"github.com/anangpermana/gin-gorm-base/controllers"
	"github.com/anangpermana/gin-gorm-base/middleware"
	"github.com/gin-gonic/gin"
)

type MemberRouteController struct {
	memberController controllers.MemberController
}

func NewMemberController(memberController controllers.MemberController) MemberRouteController {
	return MemberRouteController{memberController}
}

func (mc *MemberRouteController) MemberRoute(rg *gin.RouterGroup) {
	router := rg.Group("members")
	router.Use(middleware.DeserializeUser())
	router.POST("/create", mc.memberController.CreateMember)
	router.GET("/", mc.memberController.GetAll)
	router.GET("/:id", mc.memberController.GetOne)
	router.PUT("/:id", mc.memberController.Update)
	router.DELETE("/:id", mc.memberController.Delete)
	router.DELETE("/multiple-delete", mc.memberController.MultipleDelete)
}
