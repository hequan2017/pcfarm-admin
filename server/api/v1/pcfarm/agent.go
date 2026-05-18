package pcfarm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	pcfarmReq "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AgentApi struct{}

// Register
// @Tags    PcfarmAgent
// @Summary Ubuntu Live Agent注册
// @accept  application/json
// @Produce application/json
// @Param   data  body      pcfarmReq.AgentRegisterRequest  true  "Agent注册信息"
// @Success 200   {object}  response.Response{msg=string}
// @Router  /pcfarm/agent/register [post]
func (api *AgentApi) Register(c *gin.Context) {
	var req pcfarmReq.AgentRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pcfarmService.AgentService.Register(req); err != nil {
		global.GVA_LOG.Error("Agent注册失败", zap.Error(err))
		response.FailWithMessage("Agent注册失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("注册成功", c)
}

// Heartbeat
// @Tags    PcfarmAgent
// @Summary Ubuntu Live Agent心跳
// @accept  application/json
// @Produce application/json
// @Param   data  body      pcfarmReq.AgentHeartbeatRequest  true  "Agent心跳"
// @Success 200   {object}  response.Response{msg=string}
// @Router  /pcfarm/agent/heartbeat [post]
func (api *AgentApi) Heartbeat(c *gin.Context) {
	var req pcfarmReq.AgentHeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pcfarmService.AgentService.Heartbeat(req); err != nil {
		global.GVA_LOG.Error("Agent心跳失败", zap.Error(err))
		response.FailWithMessage("Agent心跳失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("心跳成功", c)
}
