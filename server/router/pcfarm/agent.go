package pcfarm

import "github.com/gin-gonic/gin"

type AgentRouter struct{}

func (r *AgentRouter) InitAgentRouter(Router *gin.RouterGroup) {
	agentRouter := Router.Group("pcfarm/agent")
	{
		agentRouter.POST("register", pcfarmApi.Register)
		agentRouter.POST("heartbeat", pcfarmApi.Heartbeat)
	}
}
