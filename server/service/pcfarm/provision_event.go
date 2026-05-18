package pcfarm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

type ProvisionEventService struct{}

func (ProvisionEventService) Write(serverID uint, eventType pcfarmModel.ProvisionEventType, success bool, message string) error {
	return global.GVA_DB.Create(&pcfarmModel.ProvisionEvent{
		ServerAssetID: serverID,
		Type:          eventType,
		Success:       success,
		Message:       message,
	}).Error
}
