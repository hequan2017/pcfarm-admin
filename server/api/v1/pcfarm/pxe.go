package pcfarm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	pcfarmSvc "github.com/flipped-aurora/gin-vue-admin/server/service/pcfarm"
	"github.com/gin-gonic/gin"
)

type PXEApi struct{}

// Refresh
// @Tags      PcfarmPXE
// @Summary   刷新PXE配置
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}
// @Router    /pcfarm/pxe/refresh [post]
func (api *PXEApi) Refresh(c *gin.Context) {
	response.OkWithMessage("PXE配置刷新请求已接收", c)
}

// Status
// @Tags      PcfarmPXE
// @Summary   获取PXE服务状态
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]string,msg=string}
// @Router    /pcfarm/pxe/status [get]
func (api *PXEApi) Status(c *gin.Context) {
	provider := &pcfarmSvc.LocalDnsmasqPXEProvider{}
	status, err := provider.Status()
	if err != nil {
		status = "unknown"
	}
	response.OkWithDetailed(gin.H{"status": status}, "获取成功", c)
}
