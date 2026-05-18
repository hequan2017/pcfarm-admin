package pcfarm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PXERouter struct{}

func (r *PXERouter) InitPXERouter(Router *gin.RouterGroup) {
	pxeRouter := Router.Group("pcfarm/pxe").Use(middleware.OperationRecord())
	{
		pxeRouter.POST("refresh", pcfarmApi.Refresh)
		pxeRouter.GET("status", pcfarmApi.Status)
	}
}
