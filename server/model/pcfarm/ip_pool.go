package pcfarm

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type IPAllocationStatus string

const (
	IPAllocationStatusActive   IPAllocationStatus = "active"
	IPAllocationStatusReleased IPAllocationStatus = "released"
)

// IPPool IP 地址池
type IPPool struct {
	global.GVA_MODEL
	Name      string `json:"name" form:"name" gorm:"column:name;comment:地址池名称"`
	CIDR      string `json:"cidr" form:"cidr" gorm:"column:cidr;comment:CIDR"`
	StartIP   string `json:"startIp" form:"startIp" gorm:"column:start_ip;comment:起始 IP"`
	EndIP     string `json:"endIp" form:"endIp" gorm:"column:end_ip;comment:结束 IP"`
	Gateway   string `json:"gateway" form:"gateway" gorm:"column:gateway;comment:网关"`
	DNS       string `json:"dns" form:"dns" gorm:"column:dns;comment:DNS"`
	BindIface string `json:"bindIface" form:"bindIface" gorm:"column:bind_iface;comment:绑定网卡"`
	Enabled   bool   `json:"enabled" form:"enabled" gorm:"column:enabled;comment:是否启用"`
}

func (IPPool) TableName() string {
	return "pcfarm_ip_pools"
}

// IPAllocation IP 分配记录
type IPAllocation struct {
	global.GVA_MODEL
	ServerAssetID uint               `json:"serverAssetId" form:"serverAssetId" gorm:"column:server_asset_id;comment:服务器资产ID"`
	PxeMac        string             `json:"pxeMac" form:"pxeMac" gorm:"column:pxe_mac;comment:PXE MAC 地址"`
	IP            string             `json:"ip" form:"ip" gorm:"column:ip;comment:IP 地址"`
	Status        IPAllocationStatus `json:"status" form:"status" gorm:"column:status;comment:分配状态"`
}

func (IPAllocation) TableName() string {
	return "pcfarm_ip_allocations"
}
