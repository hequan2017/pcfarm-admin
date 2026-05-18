package pcfarm

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type BootPolicy string

const (
	BootPolicyLocalDisk   BootPolicy = "local_disk"
	BootPolicyUbuntuLive  BootPolicy = "ubuntu_live"
	BootPolicyMaintenance BootPolicy = "maintenance"
)

type PowerProtocol string

const (
	PowerProtocolIPMI    PowerProtocol = "ipmi"
	PowerProtocolRedfish PowerProtocol = "redfish"
)

type ServerStatus string

const (
	ServerStatusOffline       ServerStatus = "offline"
	ServerStatusPXEReady      ServerStatus = "pxe_ready"
	ServerStatusBooting       ServerStatus = "booting"
	ServerStatusOnline        ServerStatus = "online"
	ServerStatusHeartbeatLost ServerStatus = "heartbeat_lost"
	ServerStatusPowerFailed   ServerStatus = "power_failed"
)

// ServerAsset 服务器资产
type ServerAsset struct {
	global.GVA_MODEL
	AssetCode         string        `json:"assetCode" form:"assetCode" gorm:"column:asset_code;comment:资产编号"`
	SerialNumber      string        `json:"serialNumber" form:"serialNumber" gorm:"column:serial_number;comment:序列号"`
	PxeMac            string        `json:"pxeMac" form:"pxeMac" gorm:"column:pxe_mac;comment:PXE MAC 地址"`
	FixedIP           string        `json:"fixedIp" form:"fixedIp" gorm:"column:fixed_ip;comment:固定 IP"`
	BmcAddress        string        `json:"bmcAddress" form:"bmcAddress" gorm:"column:bmc_address;comment:BMC 地址"`
	BmcUsername       string        `json:"bmcUsername" form:"bmcUsername" gorm:"column:bmc_username;comment:BMC 用户名"`
	BmcPasswordCipher string        `json:"-" form:"bmcPasswordCipher" gorm:"column:bmc_password_cipher;comment:BMC 密码密文"`
	PowerProtocol     PowerProtocol `json:"powerProtocol" form:"powerProtocol" gorm:"column:power_protocol;comment:电源协议"`
	BootPolicy        BootPolicy    `json:"bootPolicy" form:"bootPolicy" gorm:"column:boot_policy;comment:启动策略"`
	Status            ServerStatus  `json:"status" form:"status" gorm:"column:status;comment:服务器状态"`
	AgentVersion      string        `json:"agentVersion" form:"agentVersion" gorm:"column:agent_version;comment:Agent 版本"`
	HardwareSummary   string        `json:"hardwareSummary" form:"hardwareSummary" gorm:"column:hardware_summary;comment:硬件摘要"`
	LastHeartbeatAt   *time.Time    `json:"lastHeartbeatAt" form:"lastHeartbeatAt" gorm:"column:last_heartbeat_at;comment:最后心跳时间"`
}

func (ServerAsset) TableName() string {
	return "pcfarm_server_assets"
}
