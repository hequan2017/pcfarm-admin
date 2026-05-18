package pcfarm

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ServerAssetApi
	IPPoolApi
	AgentApi
	PXEApi
}

var pcfarmService = service.ServiceGroupApp.PcfarmServiceGroup
