package pcfarm

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	ServerAssetRouter
	IPPoolRouter
	AgentRouter
	PXERouter
}

var pcfarmApi = api.ApiGroupApp.PcfarmApiGroup
