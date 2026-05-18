package pcfarm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
	pcfarmReq "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IPPoolApi struct{}

// CreateIPPool
// @Tags      PcfarmIPPool
// @Summary   创建IP地址池
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      pcfarmReq.CreateIPPool  true  "IP地址池"
// @Success   200   {object}  response.Response{msg=string}
// @Router    /pcfarm/ipPool/create [post]
func (api *IPPoolApi) CreateIPPool(c *gin.Context) {
	var req pcfarmReq.CreateIPPool
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pool := pcfarmModel.IPPool{
		Name:      req.Name,
		CIDR:      req.CIDR,
		StartIP:   req.StartIP,
		EndIP:     req.EndIP,
		Gateway:   req.Gateway,
		DNS:       req.DNS,
		BindIface: req.BindIface,
		Enabled:   req.Enabled,
	}
	if err := global.GVA_DB.Create(&pool).Error; err != nil {
		global.GVA_LOG.Error("创建IP地址池失败", zap.Error(err))
		response.FailWithMessage("创建IP地址池失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetIPPoolList
// @Tags      PcfarmIPPool
// @Summary   获取IP地址池列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     pcfarmReq.IPPoolSearch  true  "分页条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /pcfarm/ipPool/list [get]
func (api *IPPoolApi) GetIPPoolList(c *gin.Context) {
	var req pcfarmReq.IPPoolSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var list []pcfarmModel.IPPool
	var total int64
	db := global.GVA_DB.Model(&pcfarmModel.IPPool{})
	if err := db.Count(&total).Error; err != nil {
		global.GVA_LOG.Error("获取IP地址池数量失败", zap.Error(err))
		response.FailWithMessage("获取IP地址池失败: "+err.Error(), c)
		return
	}
	if err := db.Order("id desc").Scopes(req.PageInfo.Paginate()).Find(&list).Error; err != nil {
		global.GVA_LOG.Error("获取IP地址池列表失败", zap.Error(err))
		response.FailWithMessage("获取IP地址池失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}
