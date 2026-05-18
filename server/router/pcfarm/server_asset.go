package pcfarm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ServerAssetRouter struct{}

func (r *ServerAssetRouter) InitServerAssetRouter(Router *gin.RouterGroup) {
	recordRouter := Router.Group("pcfarm/server").Use(middleware.OperationRecord())
	queryRouter := Router.Group("pcfarm/server")
	{
		recordRouter.POST("create", pcfarmApi.CreateServerAsset)
		recordRouter.PUT("bootPolicy", pcfarmApi.UpdateBootPolicy)
		recordRouter.POST("powerAction", pcfarmApi.ExecutePowerAction)
	}
	{
		queryRouter.GET("list", pcfarmApi.GetServerAssetList)
	}
}
