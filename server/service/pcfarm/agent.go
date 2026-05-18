package pcfarm

import (
	"errors"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
	pcfarmRequest "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/request"
	"gorm.io/gorm"
)

var ErrEmptyAgentToken = errors.New("agent token is required")

type AgentService struct{}

func (AgentService) Register(req pcfarmRequest.AgentRegisterRequest) error {
	if strings.TrimSpace(req.Token) == "" {
		return ErrEmptyAgentToken
	}

	var asset pcfarmModel.ServerAsset
	if err := findAssetBySerialOrPXEMAC(req.SerialNumber, req.PxeMac).First(&asset).Error; err != nil {
		return err
	}

	now := time.Now()
	updates := map[string]interface{}{
		"fixed_ip":          req.IP,
		"agent_version":     req.AgentVersion,
		"hardware_summary":  req.HardwareSummary,
		"last_heartbeat_at": &now,
		"status":            pcfarmModel.ServerStatusOnline,
	}
	if err := global.GVA_DB.Model(&asset).Updates(updates).Error; err != nil {
		return err
	}
	return (ProvisionEventService{}).Write(asset.ID, pcfarmModel.ProvisionEventAgentRegister, true, "agent registered")
}

func (AgentService) Heartbeat(req pcfarmRequest.AgentHeartbeatRequest) error {
	if strings.TrimSpace(req.Token) == "" {
		return ErrEmptyAgentToken
	}

	var asset pcfarmModel.ServerAsset
	if err := findAssetBySerialOrPXEMAC(req.SerialNumber, req.PxeMac).First(&asset).Error; err != nil {
		return err
	}

	now := time.Now()
	return global.GVA_DB.Model(&asset).Updates(map[string]interface{}{
		"last_heartbeat_at": &now,
		"status":            pcfarmModel.ServerStatusOnline,
	}).Error
}

func findAssetBySerialOrPXEMAC(serialNumber, pxeMAC string) *gorm.DB {
	db := global.GVA_DB.Model(&pcfarmModel.ServerAsset{})
	serialNumber = strings.TrimSpace(serialNumber)
	pxeMAC = strings.TrimSpace(pxeMAC)
	switch {
	case serialNumber != "" && pxeMAC != "":
		return db.Where("serial_number = ? OR pxe_mac = ?", serialNumber, pxeMAC)
	case serialNumber != "":
		return db.Where("serial_number = ?", serialNumber)
	default:
		return db.Where("pxe_mac = ?", pxeMAC)
	}
}
