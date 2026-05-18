package pcfarm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type IPPoolRouter struct{}

func (r *IPPoolRouter) InitIPPoolRouter(Router *gin.RouterGroup) {
	recordRouter := Router.Group("pcfarm/ipPool").Use(middleware.OperationRecord())
	queryRouter := Router.Group("pcfarm/ipPool")
	{
		recordRouter.POST("create", pcfarmApi.CreateIPPool)
	}
	{
		queryRouter.GET("list", pcfarmApi.GetIPPoolList)
	}
}
