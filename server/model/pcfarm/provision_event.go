package pcfarm

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type ProvisionEventType string

const (
	ProvisionEventAssetCreated  ProvisionEventType = "asset_created"
	ProvisionEventIPAllocated   ProvisionEventType = "ip_allocated"
	ProvisionEventPXERefreshed  ProvisionEventType = "pxe_refreshed"
	ProvisionEventPowerAction   ProvisionEventType = "power_action"
	ProvisionEventAgentRegister ProvisionEventType = "agent_register"
	ProvisionEventHeartbeatLost ProvisionEventType = "heartbeat_lost"
)

// ProvisionEvent 装机流程事件
type ProvisionEvent struct {
	global.GVA_MODEL
	ServerAssetID uint               `json:"serverAssetId" form:"serverAssetId" gorm:"column:server_asset_id;comment:服务器资产ID"`
	Type          ProvisionEventType `json:"type" form:"type" gorm:"column:type;comment:事件类型"`
	Success       bool               `json:"success" form:"success" gorm:"column:success;comment:是否成功"`
	Message       string             `json:"message" form:"message" gorm:"column:message;comment:事件消息"`
}

func (ProvisionEvent) TableName() string {
	return "pcfarm_provision_events"
}
