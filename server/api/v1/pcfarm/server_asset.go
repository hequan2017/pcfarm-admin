package pcfarm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
	pcfarmReq "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/request"
	pcfarmRes "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/response"
	pcfarmSvc "github.com/flipped-aurora/gin-vue-admin/server/service/pcfarm"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ServerAssetApi struct{}

// CreateServerAsset
// @Tags      PcfarmServer
// @Summary   创建服务器资产
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      pcfarmReq.CreateServerAsset  true  "服务器资产"
// @Success   200   {object}  response.Response{data=pcfarmRes.ServerAssetResponse,msg=string}
// @Router    /pcfarm/server/create [post]
func (api *ServerAssetApi) CreateServerAsset(c *gin.Context) {
	var req pcfarmReq.CreateServerAsset
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	asset, err := pcfarmService.ServerAssetService.CreateServerAsset(req)
	if err != nil {
		global.GVA_LOG.Error("创建服务器资产失败", zap.Error(err))
		response.FailWithMessage("创建服务器资产失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(pcfarmRes.ServerAssetResponse{Server: asset}, "创建成功", c)
}

// GetServerAssetList
// @Tags      PcfarmServer
// @Summary   获取服务器资产列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     pcfarmReq.ServerAssetSearch  true  "分页和筛选条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /pcfarm/server/list [get]
func (api *ServerAssetApi) GetServerAssetList(c *gin.Context) {
	var req pcfarmReq.ServerAssetSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pcfarmService.ServerAssetService.GetServerAssetList(req)
	if err != nil {
		global.GVA_LOG.Error("获取服务器资产列表失败", zap.Error(err))
		response.FailWithMessage("获取服务器资产列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// UpdateBootPolicy
// @Tags      PcfarmServer
// @Summary   更新服务器启动策略
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      pcfarmReq.UpdateBootPolicy  true  "启动策略"
// @Success   200   {object}  response.Response{msg=string}
// @Router    /pcfarm/server/bootPolicy [put]
func (api *ServerAssetApi) UpdateBootPolicy(c *gin.Context) {
	var req pcfarmReq.UpdateBootPolicy
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pcfarmService.ServerAssetService.UpdateBootPolicy(req.ID, pcfarmModel.BootPolicy(req.BootPolicy)); err != nil {
		global.GVA_LOG.Error("更新启动策略失败", zap.Error(err))
		response.FailWithMessage("更新启动策略失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// ExecutePowerAction
// @Tags      PcfarmServer
// @Summary   执行服务器电源动作
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      pcfarmReq.PowerActionRequest  true  "电源动作"
// @Success   200   {object}  response.Response{msg=string}
// @Router    /pcfarm/server/powerAction [post]
func (api *ServerAssetApi) ExecutePowerAction(c *gin.Context) {
	var req pcfarmReq.PowerActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pcfarmService.ServerAssetService.ExecutePowerAction(req.ID, pcfarmSvc.PowerAction(req.Action)); err != nil {
		global.GVA_LOG.Error("执行电源动作失败", zap.Error(err))
		response.FailWithMessage("执行电源动作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("执行成功", c)
}
