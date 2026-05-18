package pcfarm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
	pcfarmRequest "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/request"
	"gorm.io/gorm"
)

var (
	ErrMissingRequiredField = errors.New("missing required field")
	ErrInvalidBootPolicy    = errors.New("invalid boot policy")
)

type ServerAssetService struct{}

func (ServerAssetService) CreateServerAsset(req pcfarmRequest.CreateServerAsset) (pcfarmModel.ServerAsset, error) {
	if strings.TrimSpace(req.AssetCode) == "" || strings.TrimSpace(req.SerialNumber) == "" || strings.TrimSpace(req.PxeMac) == "" {
		return pcfarmModel.ServerAsset{}, ErrMissingRequiredField
	}

	var asset pcfarmModel.ServerAsset
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		ip, err := allocateFixedIP(tx)
		if err != nil {
			return err
		}

		asset = pcfarmModel.ServerAsset{
			AssetCode:         req.AssetCode,
			SerialNumber:      req.SerialNumber,
			PxeMac:            req.PxeMac,
			FixedIP:           ip,
			BmcAddress:        req.BmcAddress,
			BmcUsername:       req.BmcUsername,
			BmcPasswordCipher: req.BmcPassword,
			PowerProtocol:     pcfarmModel.PowerProtocol(req.PowerProtocol),
			BootPolicy:        pcfarmModel.BootPolicyLocalDisk,
			Status:            pcfarmModel.ServerStatusOffline,
		}
		if err := tx.Create(&asset).Error; err != nil {
			return err
		}
		if err := tx.Create(&pcfarmModel.IPAllocation{
			ServerAssetID: asset.ID,
			PxeMac:        asset.PxeMac,
			IP:            ip,
			Status:        pcfarmModel.IPAllocationStatusActive,
		}).Error; err != nil {
			return err
		}
		return tx.Create(&pcfarmModel.ProvisionEvent{
			ServerAssetID: asset.ID,
			Type:          pcfarmModel.ProvisionEventAssetCreated,
			Success:       true,
			Message:       "asset created",
		}).Error
	})
	return asset, err
}

func (ServerAssetService) GetServerAssetList(info pcfarmRequest.ServerAssetSearch) (list []pcfarmModel.ServerAsset, total int64, err error) {
	db := global.GVA_DB.Model(&pcfarmModel.ServerAsset{})
	keyword := strings.TrimSpace(info.Keyword)
	if keyword == "" {
		keyword = strings.TrimSpace(info.PageInfo.Keyword)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("asset_code LIKE ? OR serial_number LIKE ? OR pxe_mac LIKE ? OR fixed_ip LIKE ?", like, like, like, like)
	}
	if strings.TrimSpace(info.Status) != "" {
		db = db.Where("status = ?", info.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("id desc").Scopes(info.PageInfo.Paginate()).Find(&list).Error
	return list, total, err
}

func (ServerAssetService) UpdateBootPolicy(id uint, policy pcfarmModel.BootPolicy) error {
	if !validBootPolicy(policy) {
		return ErrInvalidBootPolicy
	}

	var asset pcfarmModel.ServerAsset
	if err := global.GVA_DB.First(&asset, id).Error; err != nil {
		return err
	}
	asset.BootPolicy = policy
	asset.Status = pcfarmModel.ServerStatusPXEReady
	if err := global.GVA_DB.Save(&asset).Error; err != nil {
		return err
	}
	if err := (LocalDnsmasqPXEProvider{}).Refresh(asset); err != nil {
		_ = (ProvisionEventService{}).Write(asset.ID, pcfarmModel.ProvisionEventPXERefreshed, false, err.Error())
		return err
	}
	return (ProvisionEventService{}).Write(asset.ID, pcfarmModel.ProvisionEventPXERefreshed, true, "pxe refreshed")
}

func (ServerAssetService) ExecutePowerAction(id uint, action PowerAction) error {
	var asset pcfarmModel.ServerAsset
	if err := global.GVA_DB.First(&asset, id).Error; err != nil {
		return err
	}
	provider, err := powerProviderFor(asset.PowerProtocol)
	if err != nil {
		_ = (ProvisionEventService{}).Write(asset.ID, pcfarmModel.ProvisionEventPowerAction, false, err.Error())
		return err
	}
	if err := provider.Execute(asset, action); err != nil {
		_ = (ProvisionEventService{}).Write(asset.ID, pcfarmModel.ProvisionEventPowerAction, false, err.Error())
		return err
	}
	return (ProvisionEventService{}).Write(asset.ID, pcfarmModel.ProvisionEventPowerAction, true, fmt.Sprintf("power action %s executed", action))
}

func allocateFixedIP(tx *gorm.DB) (string, error) {
	var pool pcfarmModel.IPPool
	if err := tx.Where("enabled = ?", true).Order("id asc").First(&pool).Error; err != nil {
		return "", err
	}

	var allocations []pcfarmModel.IPAllocation
	if err := tx.Where("status = ?", pcfarmModel.IPAllocationStatusActive).Find(&allocations).Error; err != nil {
		return "", err
	}
	allocated := make(map[string]struct{}, len(allocations))
	for _, allocation := range allocations {
		allocated[allocation.IP] = struct{}{}
	}
	return nextAvailableIP(pool, allocated)
}

func validBootPolicy(policy pcfarmModel.BootPolicy) bool {
	switch policy {
	case pcfarmModel.BootPolicyLocalDisk, pcfarmModel.BootPolicyUbuntuLive, pcfarmModel.BootPolicyMaintenance:
		return true
	default:
		return false
	}
}
