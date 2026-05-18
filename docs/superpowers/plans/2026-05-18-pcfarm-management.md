# pcfarm Management Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a pcfarm management module that manages server assets, fixed IP allocation, PXE boot policy, IPMI/Redfish power actions, PXE configuration generation, and Ubuntu Live agent registration.

**Architecture:** Implement `pcfarm` as a first-class business module following the existing `Router -> API -> Service -> Model` structure. Keep PXE, power control, and IP allocation behind small provider interfaces so the MVP can use local `dnsmasq` and command adapters while preserving a clean path to Kea DHCP or a remote PXE agent later.

**Tech Stack:** Go, Gin, GORM, Swagger comments, Vue 3, Vite, Element Plus, existing gin-vue-admin response/request utilities.

---

## Scope And Rules

- Do not create git commits. The project rules explicitly say not to plan or execute commits unless the user asks.
- Keep comments in Chinese where new explanatory comments are needed.
- Do not implement multi-image Ubuntu management, multi-PXE-node orchestration, DHCP relay, or local-disk reinstall.
- Do not read or modify `node_modules/`.
- Prefer unit tests for business rules before wiring UI.

## File Structure

Create backend model files:

- `server/model/pcfarm/server_asset.go`: server asset, boot policy, power protocol, and server status models.
- `server/model/pcfarm/ip_pool.go`: IP pool and IP allocation models.
- `server/model/pcfarm/provision_event.go`: audit and lifecycle event model.
- `server/model/pcfarm/request/server_asset.go`: server asset create/update/search/action request structs.
- `server/model/pcfarm/request/ip_pool.go`: IP pool create/update/search request structs.
- `server/model/pcfarm/request/agent.go`: Ubuntu Live agent register and heartbeat request structs.
- `server/model/pcfarm/response/server_asset.go`: server detail/list response structs with masked credentials.
- `server/model/pcfarm/response/ip_pool.go`: IP pool summary response structs.

Create backend service files:

- `server/service/pcfarm/enter.go`: service group registration.
- `server/service/pcfarm/ip_allocator.go`: deterministic fixed IP allocation.
- `server/service/pcfarm/server_asset.go`: asset CRUD and boot policy workflow.
- `server/service/pcfarm/pxe_provider.go`: PXE provider interface and local dnsmasq implementation.
- `server/service/pcfarm/power_provider.go`: power provider interface with IPMI and Redfish adapters.
- `server/service/pcfarm/agent.go`: Ubuntu Live agent registration and heartbeat handling.
- `server/service/pcfarm/provision_event.go`: event writing helpers.

Create backend API and router files:

- `server/api/v1/pcfarm/enter.go`: API group registration.
- `server/api/v1/pcfarm/server_asset.go`: asset, boot policy, and power action HTTP handlers.
- `server/api/v1/pcfarm/ip_pool.go`: IP pool HTTP handlers.
- `server/api/v1/pcfarm/agent.go`: Live agent HTTP handlers.
- `server/api/v1/pcfarm/pxe.go`: PXE refresh and service status handlers.
- `server/router/pcfarm/enter.go`: router group registration.
- `server/router/pcfarm/server_asset.go`: asset route setup.
- `server/router/pcfarm/ip_pool.go`: IP pool route setup.
- `server/router/pcfarm/agent.go`: agent route setup.
- `server/router/pcfarm/pxe.go`: PXE route setup.

Modify backend integration files:

- `server/model/pcfarm/enter.go`: optional model list helper if consistent with implementation.
- `server/service/enter.go`: expose `PcfarmServiceGroupApp`.
- `server/api/v1/enter.go`: expose `PcfarmApiGroupApp`.
- `server/router/enter.go`: expose `PcfarmRouterGroupApp`.
- `server/initialize/router.go`: register pcfarm routes.
- `server/initialize/gorm.go` or `server/initialize/ensure_tables.go`: auto-migrate pcfarm tables using the project’s existing pattern.

Create frontend files:

- `web/src/api/pcfarm.js`: pcfarm API wrapper functions.
- `web/src/view/pcfarm/server/index.vue`: server asset list and batch actions.
- `web/src/view/pcfarm/server/detail.vue`: server detail, credentials reset, event timeline.
- `web/src/view/pcfarm/ipPool/index.vue`: IP pool management.
- `web/src/view/pcfarm/pxe/index.vue`: PXE settings and refresh/status controls.

Modify frontend routing/menu data only through the project’s existing menu initialization or admin menu mechanism. If menus are DB-driven in this repo, add a follow-up SQL or initialization note instead of hard-coding routes in unrelated files.

## Task 1: Backend Domain Models

**Files:**

- Create: `server/model/pcfarm/server_asset.go`
- Create: `server/model/pcfarm/ip_pool.go`
- Create: `server/model/pcfarm/provision_event.go`
- Create: `server/model/pcfarm/request/server_asset.go`
- Create: `server/model/pcfarm/request/ip_pool.go`
- Create: `server/model/pcfarm/request/agent.go`
- Create: `server/model/pcfarm/response/server_asset.go`
- Create: `server/model/pcfarm/response/ip_pool.go`

- [ ] **Step 1: Create pcfarm domain enums and asset model**

Add `server/model/pcfarm/server_asset.go`:

```go
package pcfarm

import "github.com/flipped-aurora/gin-vue-admin/server/global"

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
	ServerStatusOffline      ServerStatus = "offline"
	ServerStatusPXEReady     ServerStatus = "pxe_ready"
	ServerStatusBooting      ServerStatus = "booting"
	ServerStatusOnline       ServerStatus = "online"
	ServerStatusHeartbeatLost ServerStatus = "heartbeat_lost"
	ServerStatusPowerFailed  ServerStatus = "power_failed"
)

type ServerAsset struct {
	global.GVA_MODEL
	AssetCode          string        `json:"assetCode" form:"assetCode" gorm:"column:asset_code;uniqueIndex;comment:资产编号"`
	SerialNumber       string        `json:"serialNumber" form:"serialNumber" gorm:"column:serial_number;uniqueIndex;comment:序列号"`
	PxeMac             string        `json:"pxeMac" form:"pxeMac" gorm:"column:pxe_mac;uniqueIndex;comment:PXE网卡MAC"`
	FixedIP            string        `json:"fixedIp" form:"fixedIp" gorm:"column:fixed_ip;comment:固定IP"`
	BmcAddress         string        `json:"bmcAddress" form:"bmcAddress" gorm:"column:bmc_address;comment:BMC地址"`
	BmcUsername        string        `json:"bmcUsername" form:"bmcUsername" gorm:"column:bmc_username;comment:BMC用户名"`
	BmcPasswordCipher  string        `json:"-" form:"-" gorm:"column:bmc_password_cipher;comment:BMC密码密文"`
	PowerProtocol      PowerProtocol `json:"powerProtocol" form:"powerProtocol" gorm:"column:power_protocol;comment:远控协议"`
	BootPolicy         BootPolicy    `json:"bootPolicy" form:"bootPolicy" gorm:"column:boot_policy;comment:启动策略"`
	Status             ServerStatus  `json:"status" form:"status" gorm:"column:status;comment:服务器状态"`
	AgentVersion       string        `json:"agentVersion" form:"agentVersion" gorm:"column:agent_version;comment:Agent版本"`
	HardwareSummary    string        `json:"hardwareSummary" form:"hardwareSummary" gorm:"column:hardware_summary;type:text;comment:硬件摘要"`
	LastHeartbeatAt    *global.LocalTime `json:"lastHeartbeatAt" form:"lastHeartbeatAt" gorm:"column:last_heartbeat_at;comment:最后心跳时间"`
}

func (ServerAsset) TableName() string {
	return "pcfarm_server_assets"
}
```

- [ ] **Step 2: Create IP pool and allocation models**

Add `server/model/pcfarm/ip_pool.go`:

```go
package pcfarm

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type IPPool struct {
	global.GVA_MODEL
	Name        string `json:"name" form:"name" gorm:"column:name;comment:地址池名称"`
	CIDR        string `json:"cidr" form:"cidr" gorm:"column:cidr;comment:网段CIDR"`
	StartIP     string `json:"startIp" form:"startIp" gorm:"column:start_ip;comment:起始IP"`
	EndIP       string `json:"endIp" form:"endIp" gorm:"column:end_ip;comment:结束IP"`
	Gateway     string `json:"gateway" form:"gateway" gorm:"column:gateway;comment:网关"`
	DNS         string `json:"dns" form:"dns" gorm:"column:dns;comment:DNS"`
	BindIface   string `json:"bindIface" form:"bindIface" gorm:"column:bind_iface;comment:绑定网卡"`
	Enabled     bool   `json:"enabled" form:"enabled" gorm:"column:enabled;comment:是否启用"`
}

func (IPPool) TableName() string {
	return "pcfarm_ip_pools"
}

type IPAllocationStatus string

const (
	IPAllocationActive   IPAllocationStatus = "active"
	IPAllocationReleased IPAllocationStatus = "released"
)

type IPAllocation struct {
	global.GVA_MODEL
	ServerAssetID uint               `json:"serverAssetId" form:"serverAssetId" gorm:"column:server_asset_id;index;comment:服务器ID"`
	PxeMac        string             `json:"pxeMac" form:"pxeMac" gorm:"column:pxe_mac;uniqueIndex:idx_pcfarm_active_mac;comment:PXE网卡MAC"`
	IP            string             `json:"ip" form:"ip" gorm:"column:ip;uniqueIndex:idx_pcfarm_active_ip;comment:分配IP"`
	Status        IPAllocationStatus `json:"status" form:"status" gorm:"column:status;index;comment:分配状态"`
}

func (IPAllocation) TableName() string {
	return "pcfarm_ip_allocations"
}
```

- [ ] **Step 3: Create event model**

Add `server/model/pcfarm/provision_event.go`:

```go
package pcfarm

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type ProvisionEventType string

const (
	ProvisionEventAssetCreated ProvisionEventType = "asset_created"
	ProvisionEventIPAllocated  ProvisionEventType = "ip_allocated"
	ProvisionEventPXERefreshed ProvisionEventType = "pxe_refreshed"
	ProvisionEventPowerAction  ProvisionEventType = "power_action"
	ProvisionEventAgentRegister ProvisionEventType = "agent_register"
	ProvisionEventHeartbeatLost ProvisionEventType = "heartbeat_lost"
)

type ProvisionEvent struct {
	global.GVA_MODEL
	ServerAssetID uint               `json:"serverAssetId" form:"serverAssetId" gorm:"column:server_asset_id;index;comment:服务器ID"`
	Type          ProvisionEventType `json:"type" form:"type" gorm:"column:type;comment:事件类型"`
	Success       bool               `json:"success" form:"success" gorm:"column:success;comment:是否成功"`
	Message       string             `json:"message" form:"message" gorm:"column:message;type:text;comment:事件信息"`
}

func (ProvisionEvent) TableName() string {
	return "pcfarm_provision_events"
}
```

- [ ] **Step 4: Create request and response structs**

Add request structs with exact names used by later tasks:

```go
package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type ServerAssetSearch struct {
	Keyword string `json:"keyword" form:"keyword"`
	Status  string `json:"status" form:"status"`
	request.PageInfo
}

type CreateServerAsset struct {
	AssetCode     string `json:"assetCode"`
	SerialNumber  string `json:"serialNumber"`
	PxeMac        string `json:"pxeMac"`
	BmcAddress    string `json:"bmcAddress"`
	BmcUsername   string `json:"bmcUsername"`
	BmcPassword   string `json:"bmcPassword"`
	PowerProtocol string `json:"powerProtocol"`
}

type UpdateBootPolicy struct {
	ID         uint   `json:"id"`
	BootPolicy string `json:"bootPolicy"`
}

type PowerActionRequest struct {
	ID     uint   `json:"id"`
	Action string `json:"action"`
}
```

Add IP pool request structs:

```go
package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type IPPoolSearch struct {
	request.PageInfo
}

type CreateIPPool struct {
	Name      string `json:"name"`
	CIDR      string `json:"cidr"`
	StartIP   string `json:"startIp"`
	EndIP     string `json:"endIp"`
	Gateway   string `json:"gateway"`
	DNS       string `json:"dns"`
	BindIface string `json:"bindIface"`
	Enabled   bool   `json:"enabled"`
}
```

Add agent request structs:

```go
package request

type AgentRegisterRequest struct {
	SerialNumber    string `json:"serialNumber"`
	PxeMac          string `json:"pxeMac"`
	IP              string `json:"ip"`
	AgentVersion    string `json:"agentVersion"`
	HardwareSummary string `json:"hardwareSummary"`
	Token           string `json:"token"`
}

type AgentHeartbeatRequest struct {
	SerialNumber string `json:"serialNumber"`
	PxeMac       string `json:"pxeMac"`
	Token        string `json:"token"`
}
```

Add response structs:

```go
package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"

type ServerAssetResponse struct {
	Server pcfarm.ServerAsset `json:"server"`
}

type ServerAssetListItem struct {
	pcfarm.ServerAsset
	HasBmcPassword bool `json:"hasBmcPassword"`
}
```

```go
package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"

type IPPoolSummary struct {
	Pool           pcfarm.IPPool `json:"pool"`
	AllocatedCount int64         `json:"allocatedCount"`
	AvailableCount int64         `json:"availableCount"`
}
```

- [ ] **Step 5: Run backend compile for model-only changes**

Run: `go test ./model/...`

Expected: package compiles. If `global.LocalTime` is not available in this repo, replace `*global.LocalTime` with `*time.Time` and import `time`.

## Task 2: IP Allocator

**Files:**

- Create: `server/service/pcfarm/ip_allocator.go`
- Test: `server/service/pcfarm/ip_allocator_test.go`

- [ ] **Step 1: Write allocator tests**

Add `server/service/pcfarm/ip_allocator_test.go`:

```go
package pcfarm

import (
	"testing"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

func TestNextAvailableIPSkipsAllocated(t *testing.T) {
	pool := model.IPPool{StartIP: "192.168.50.10", EndIP: "192.168.50.12"}
	allocated := map[string]struct{}{"192.168.50.10": {}}

	got, err := nextAvailableIP(pool, allocated)
	if err != nil {
		t.Fatalf("nextAvailableIP returned error: %v", err)
	}
	if got != "192.168.50.11" {
		t.Fatalf("expected 192.168.50.11, got %s", got)
	}
}

func TestNextAvailableIPReturnsErrorWhenPoolExhausted(t *testing.T) {
	pool := model.IPPool{StartIP: "192.168.50.10", EndIP: "192.168.50.11"}
	allocated := map[string]struct{}{
		"192.168.50.10": {},
		"192.168.50.11": {},
	}

	_, err := nextAvailableIP(pool, allocated)
	if err == nil {
		t.Fatal("expected pool exhausted error")
	}
}
```

- [ ] **Step 2: Run tests and verify failure**

Run: `go test ./service/pcfarm -run TestNextAvailableIP`

Expected: FAIL because `nextAvailableIP` is not defined.

- [ ] **Step 3: Implement pure IP selection helper**

Add `server/service/pcfarm/ip_allocator.go`:

```go
package pcfarm

import (
	"encoding/binary"
	"errors"
	"net"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

var ErrIPPoolExhausted = errors.New("IP地址池已耗尽")

func nextAvailableIP(pool model.IPPool, allocated map[string]struct{}) (string, error) {
	start := net.ParseIP(pool.StartIP).To4()
	end := net.ParseIP(pool.EndIP).To4()
	if start == nil || end == nil {
		return "", errors.New("地址池IP格式不正确")
	}

	startNum := binary.BigEndian.Uint32(start)
	endNum := binary.BigEndian.Uint32(end)
	if startNum > endNum {
		return "", errors.New("地址池起始IP不能大于结束IP")
	}

	for current := startNum; current <= endNum; current++ {
		var raw [4]byte
		binary.BigEndian.PutUint32(raw[:], current)
		ip := net.IP(raw[:]).String()
		if _, exists := allocated[ip]; !exists {
			return ip, nil
		}
	}

	return "", ErrIPPoolExhausted
}
```

- [ ] **Step 4: Fix the compile issue in the implementation**

If Go reports that `binary.BigEndian.PutUint32` cannot accept `raw[:]`, adjust to:

```go
binary.BigEndian.PutUint32(raw[:], current)
```

Expected final code compiles with `raw[:]` because it is a `[]byte` view over `[4]byte`.

- [ ] **Step 5: Run allocator tests**

Run: `go test ./service/pcfarm -run TestNextAvailableIP`

Expected: PASS.

## Task 3: Asset Service And IP Allocation Workflow

**Files:**

- Modify: `server/service/pcfarm/ip_allocator.go`
- Create: `server/service/pcfarm/provision_event.go`
- Create: `server/service/pcfarm/server_asset.go`
- Create: `server/service/pcfarm/enter.go`
- Test: `server/service/pcfarm/server_asset_test.go`

- [ ] **Step 1: Add service group**

Create `server/service/pcfarm/enter.go`:

```go
package pcfarm

type ServiceGroup struct {
	ServerAssetService
	AgentService
	ProvisionEventService
}
```

- [ ] **Step 2: Add event helper**

Create `server/service/pcfarm/provision_event.go`:

```go
package pcfarm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

type ProvisionEventService struct{}

func (s *ProvisionEventService) Write(serverID uint, eventType model.ProvisionEventType, success bool, message string) error {
	return global.GVA_DB.Create(&model.ProvisionEvent{
		ServerAssetID: serverID,
		Type:          eventType,
		Success:       success,
		Message:       message,
	}).Error
}
```

- [ ] **Step 3: Implement asset create workflow**

Create `server/service/pcfarm/server_asset.go`:

```go
package pcfarm

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
	pcfarmReq "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/request"
	"gorm.io/gorm"
)

type ServerAssetService struct{}

func (s *ServerAssetService) CreateServerAsset(req pcfarmReq.CreateServerAsset) (model.ServerAsset, error) {
	if req.AssetCode == "" || req.SerialNumber == "" || req.PxeMac == "" {
		return model.ServerAsset{}, errors.New("资产编号、序列号和PXE MAC不能为空")
	}

	var asset model.ServerAsset
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		ip, err := allocateFixedIP(tx, req.PxeMac)
		if err != nil {
			return err
		}

		asset = model.ServerAsset{
			AssetCode:         req.AssetCode,
			SerialNumber:      req.SerialNumber,
			PxeMac:            req.PxeMac,
			FixedIP:           ip,
			BmcAddress:        req.BmcAddress,
			BmcUsername:       req.BmcUsername,
			BmcPasswordCipher: req.BmcPassword,
			PowerProtocol:     model.PowerProtocol(req.PowerProtocol),
			BootPolicy:        model.BootPolicyLocalDisk,
			Status:            model.ServerStatusOffline,
		}
		if err := tx.Create(&asset).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.IPAllocation{
			ServerAssetID: asset.ID,
			PxeMac:        req.PxeMac,
			IP:            ip,
			Status:        model.IPAllocationActive,
		}).Error; err != nil {
			return err
		}

		return tx.Create(&model.ProvisionEvent{
			ServerAssetID: asset.ID,
			Type:          model.ProvisionEventAssetCreated,
			Success:       true,
			Message:       "服务器资产已创建并完成固定IP分配",
		}).Error
	})

	return asset, err
}

func allocateFixedIP(tx *gorm.DB, pxeMac string) (string, error) {
	var exists model.IPAllocation
	if err := tx.Where("pxe_mac = ? AND status = ?", pxeMac, model.IPAllocationActive).First(&exists).Error; err == nil {
		return exists.IP, nil
	}

	var pool model.IPPool
	if err := tx.Where("enabled = ?", true).First(&pool).Error; err != nil {
		return "", errors.New("未配置启用的IP地址池")
	}

	var allocations []model.IPAllocation
	if err := tx.Where("status = ?", model.IPAllocationActive).Find(&allocations).Error; err != nil {
		return "", err
	}

	allocated := make(map[string]struct{}, len(allocations))
	for _, allocation := range allocations {
		allocated[allocation.IP] = struct{}{}
	}

	return nextAvailableIP(pool, allocated)
}
```

- [ ] **Step 4: Add list and boot policy service methods**

Append to `server/service/pcfarm/server_asset.go`:

```go
func (s *ServerAssetService) GetServerAssetList(info pcfarmReq.ServerAssetSearch) (list []model.ServerAsset, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&model.ServerAsset{})

	if info.Keyword != "" {
		like := "%" + info.Keyword + "%"
		db = db.Where("asset_code LIKE ? OR serial_number LIKE ? OR pxe_mac LIKE ? OR fixed_ip LIKE ?", like, like, like, like)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}

	err = db.Limit(limit).Offset(offset).Order("id DESC").Find(&list).Error
	return list, total, err
}

func (s *ServerAssetService) UpdateBootPolicy(id uint, policy model.BootPolicy) error {
	if policy != model.BootPolicyLocalDisk && policy != model.BootPolicyUbuntuLive && policy != model.BootPolicyMaintenance {
		return errors.New("不支持的启动策略")
	}
	return global.GVA_DB.Model(&model.ServerAsset{}).Where("id = ?", id).Updates(map[string]interface{}{
		"boot_policy": policy,
		"status":      model.ServerStatusPXEReady,
	}).Error
}
```

- [ ] **Step 5: Run focused service tests**

Run: `go test ./service/pcfarm`

Expected: PASS after adding any required imports and fixing formatting with `gofmt`.

## Task 4: PXE Provider

**Files:**

- Create: `server/service/pcfarm/pxe_provider.go`
- Test: `server/service/pcfarm/pxe_provider_test.go`

- [ ] **Step 1: Write PXE config rendering tests**

Create `server/service/pcfarm/pxe_provider_test.go`:

```go
package pcfarm

import (
	"strings"
	"testing"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

func TestRenderDnsmasqHostBinding(t *testing.T) {
	server := model.ServerAsset{PxeMac: "52:54:00:12:34:56", FixedIP: "192.168.50.20"}

	got := renderDnsmasqHostBinding(server)
	if !strings.Contains(got, "52:54:00:12:34:56,192.168.50.20") {
		t.Fatalf("dnsmasq binding missing MAC/IP: %s", got)
	}
}

func TestRenderBootMenuUsesUbuntuLivePolicy(t *testing.T) {
	server := model.ServerAsset{BootPolicy: model.BootPolicyUbuntuLive}

	got := renderBootMenu(server)
	if !strings.Contains(got, "Ubuntu Live") {
		t.Fatalf("expected Ubuntu Live menu, got %s", got)
	}
}
```

- [ ] **Step 2: Implement provider interface and local renderer**

Create `server/service/pcfarm/pxe_provider.go`:

```go
package pcfarm

import (
	"fmt"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

type PXEProvider interface {
	Refresh(server model.ServerAsset) error
	Status() (string, error)
}

type LocalDnsmasqPXEProvider struct {
	ConfigDir string
}

func (p *LocalDnsmasqPXEProvider) Refresh(server model.ServerAsset) error {
	_ = renderDnsmasqHostBinding(server)
	_ = renderBootMenu(server)
	return nil
}

func (p *LocalDnsmasqPXEProvider) Status() (string, error) {
	return "unknown", nil
}

func renderDnsmasqHostBinding(server model.ServerAsset) string {
	return fmt.Sprintf("dhcp-host=%s,%s\n", server.PxeMac, server.FixedIP)
}

func renderBootMenu(server model.ServerAsset) string {
	switch server.BootPolicy {
	case model.BootPolicyUbuntuLive:
		return "menuentry 'Ubuntu Live' { linux /ubuntu/vmlinuz boot=casper netboot=http }\n"
	case model.BootPolicyMaintenance:
		return "menuentry 'Maintenance' { linux /ubuntu/vmlinuz boot=casper pcfarm.mode=maintenance }\n"
	default:
		return "menuentry 'Local Disk' { exit }\n"
	}
}
```

- [ ] **Step 3: Run PXE provider tests**

Run: `go test ./service/pcfarm -run TestRender`

Expected: PASS.

- [ ] **Step 4: Wire boot policy update to PXE refresh**

Modify `UpdateBootPolicy` in `server/service/pcfarm/server_asset.go` so it loads the server after update and calls `LocalDnsmasqPXEProvider.Refresh(server)`.

Use this exact method body:

```go
func (s *ServerAssetService) UpdateBootPolicy(id uint, policy model.BootPolicy) error {
	if policy != model.BootPolicyLocalDisk && policy != model.BootPolicyUbuntuLive && policy != model.BootPolicyMaintenance {
		return errors.New("不支持的启动策略")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.ServerAsset{}).Where("id = ?", id).Updates(map[string]interface{}{
			"boot_policy": policy,
			"status":      model.ServerStatusPXEReady,
		}).Error; err != nil {
			return err
		}

		var server model.ServerAsset
		if err := tx.Where("id = ?", id).First(&server).Error; err != nil {
			return err
		}

		provider := &LocalDnsmasqPXEProvider{}
		if err := provider.Refresh(server); err != nil {
			return err
		}

		return tx.Create(&model.ProvisionEvent{
			ServerAssetID: id,
			Type:          model.ProvisionEventPXERefreshed,
			Success:       true,
			Message:       "PXE配置已刷新",
		}).Error
	})
}
```

Run: `go test ./service/pcfarm`

Expected: PASS.

## Task 5: Power Provider

**Files:**

- Create: `server/service/pcfarm/power_provider.go`
- Test: `server/service/pcfarm/power_provider_test.go`

- [ ] **Step 1: Write provider selection tests**

Create `server/service/pcfarm/power_provider_test.go`:

```go
package pcfarm

import (
	"testing"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

func TestPowerProviderForIPMI(t *testing.T) {
	provider, err := powerProviderFor(model.PowerProtocolIPMI)
	if err != nil {
		t.Fatalf("powerProviderFor returned error: %v", err)
	}
	if _, ok := provider.(*IPMIPowerProvider); !ok {
		t.Fatalf("expected IPMIPowerProvider, got %T", provider)
	}
}

func TestPowerProviderForUnsupportedProtocol(t *testing.T) {
	_, err := powerProviderFor("unknown")
	if err == nil {
		t.Fatal("expected unsupported protocol error")
	}
}
```

- [ ] **Step 2: Implement power provider abstraction**

Create `server/service/pcfarm/power_provider.go`:

```go
package pcfarm

import (
	"errors"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

type PowerAction string

const (
	PowerActionOn      PowerAction = "on"
	PowerActionOff     PowerAction = "off"
	PowerActionReboot  PowerAction = "reboot"
	PowerActionBootPXE PowerAction = "boot_pxe"
)

type PowerProvider interface {
	Execute(server model.ServerAsset, action PowerAction) error
}

type IPMIPowerProvider struct{}

func (p *IPMIPowerProvider) Execute(server model.ServerAsset, action PowerAction) error {
	return nil
}

type RedfishPowerProvider struct{}

func (p *RedfishPowerProvider) Execute(server model.ServerAsset, action PowerAction) error {
	return nil
}

func powerProviderFor(protocol model.PowerProtocol) (PowerProvider, error) {
	switch protocol {
	case model.PowerProtocolIPMI:
		return &IPMIPowerProvider{}, nil
	case model.PowerProtocolRedfish:
		return &RedfishPowerProvider{}, nil
	default:
		return nil, errors.New("不支持的远控协议")
	}
}
```

- [ ] **Step 3: Add service method for power action**

Append to `server/service/pcfarm/server_asset.go`:

```go
func (s *ServerAssetService) ExecutePowerAction(id uint, action PowerAction) error {
	var server model.ServerAsset
	if err := global.GVA_DB.Where("id = ?", id).First(&server).Error; err != nil {
		return err
	}

	provider, err := powerProviderFor(server.PowerProtocol)
	if err != nil {
		return err
	}

	if err := provider.Execute(server, action); err != nil {
		_ = global.GVA_DB.Create(&model.ProvisionEvent{
			ServerAssetID: id,
			Type:          model.ProvisionEventPowerAction,
			Success:       false,
			Message:       err.Error(),
		}).Error
		return err
	}

	return global.GVA_DB.Create(&model.ProvisionEvent{
		ServerAssetID: id,
		Type:          model.ProvisionEventPowerAction,
		Success:       true,
		Message:       string(action),
	}).Error
}
```

- [ ] **Step 4: Run power tests**

Run: `go test ./service/pcfarm -run TestPowerProvider`

Expected: PASS.

## Task 6: Agent Registration

**Files:**

- Create: `server/service/pcfarm/agent.go`
- Test: `server/service/pcfarm/agent_test.go`

- [ ] **Step 1: Implement agent service**

Create `server/service/pcfarm/agent.go`:

```go
package pcfarm

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
	pcfarmReq "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/request"
)

type AgentService struct{}

func (s *AgentService) Register(req pcfarmReq.AgentRegisterRequest) error {
	if req.Token == "" {
		return errors.New("agent token不能为空")
	}

	var server model.ServerAsset
	if err := global.GVA_DB.Where("serial_number = ? OR pxe_mac = ?", req.SerialNumber, req.PxeMac).First(&server).Error; err != nil {
		return err
	}

	now := time.Now()
	if err := global.GVA_DB.Model(&server).Updates(map[string]interface{}{
		"fixed_ip":           req.IP,
		"agent_version":      req.AgentVersion,
		"hardware_summary":   req.HardwareSummary,
		"last_heartbeat_at":  &now,
		"status":             model.ServerStatusOnline,
	}).Error; err != nil {
		return err
	}

	return global.GVA_DB.Create(&model.ProvisionEvent{
		ServerAssetID: server.ID,
		Type:          model.ProvisionEventAgentRegister,
		Success:       true,
		Message:       "Ubuntu Live Agent已注册",
	}).Error
}

func (s *AgentService) Heartbeat(req pcfarmReq.AgentHeartbeatRequest) error {
	if req.Token == "" {
		return errors.New("agent token不能为空")
	}

	now := time.Now()
	return global.GVA_DB.Model(&model.ServerAsset{}).
		Where("serial_number = ? OR pxe_mac = ?", req.SerialNumber, req.PxeMac).
		Updates(map[string]interface{}{
			"last_heartbeat_at": &now,
			"status":            model.ServerStatusOnline,
		}).Error
}
```

- [ ] **Step 2: Run agent service compile**

Run: `go test ./service/pcfarm -run TestDoesNotExist`

Expected: PASS with no tests to run, confirming package compiles.

## Task 7: Backend API And Router

**Files:**

- Create: `server/api/v1/pcfarm/enter.go`
- Create: `server/api/v1/pcfarm/server_asset.go`
- Create: `server/api/v1/pcfarm/ip_pool.go`
- Create: `server/api/v1/pcfarm/agent.go`
- Create: `server/router/pcfarm/enter.go`
- Create: `server/router/pcfarm/server_asset.go`
- Create: `server/router/pcfarm/ip_pool.go`
- Create: `server/router/pcfarm/agent.go`
- Modify: `server/api/v1/enter.go`
- Modify: `server/router/enter.go`
- Modify: `server/initialize/router.go`

- [ ] **Step 1: Create pcfarm API group**

Create `server/api/v1/pcfarm/enter.go`:

```go
package pcfarm

import pcfarmService "github.com/flipped-aurora/gin-vue-admin/server/service/pcfarm"

type ApiGroup struct {
	ServerAssetApi
	IPPoolApi
	AgentApi
}

var pcfarmServiceGroup = pcfarmService.ServiceGroup{}
```

- [ ] **Step 2: Create server asset API**

Create `server/api/v1/pcfarm/server_asset.go` with handlers for create, list, boot policy, and power action. Use `response.OkWithDetailed` for list results and `response.FailWithMessage` for errors.

Core create handler:

```go
func (api *ServerAssetApi) CreateServerAsset(c *gin.Context) {
	var req pcfarmReq.CreateServerAsset
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	asset, err := pcfarmServiceGroup.ServerAssetService.CreateServerAsset(req)
	if err != nil {
		global.GVA_LOG.Error("创建服务器资产失败", zap.Error(err))
		response.FailWithMessage("创建服务器资产失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(pcfarmRes.ServerAssetResponse{Server: asset}, "创建成功", c)
}
```

Core list handler:

```go
func (api *ServerAssetApi) GetServerAssetList(c *gin.Context) {
	var req pcfarmReq.ServerAssetSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pcfarmServiceGroup.ServerAssetService.GetServerAssetList(req)
	if err != nil {
		response.FailWithMessage("获取服务器列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}
```

- [ ] **Step 3: Create agent API**

Create `server/api/v1/pcfarm/agent.go`:

```go
package pcfarm

import (
	pcfarmReq "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type AgentApi struct{}

func (api *AgentApi) Register(c *gin.Context) {
	var req pcfarmReq.AgentRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pcfarmServiceGroup.AgentService.Register(req); err != nil {
		response.FailWithMessage("Agent注册失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("注册成功", c)
}

func (api *AgentApi) Heartbeat(c *gin.Context) {
	var req pcfarmReq.AgentHeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pcfarmServiceGroup.AgentService.Heartbeat(req); err != nil {
		response.FailWithMessage("心跳失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("心跳成功", c)
}
```

- [ ] **Step 4: Create routers and register them**

Create route files using existing project style:

```go
package pcfarm

type RouterGroup struct {
	ServerAssetRouter
	IPPoolRouter
	AgentRouter
}
```

Server routes:

```go
func (r *ServerAssetRouter) InitServerAssetRouter(Router *gin.RouterGroup) {
	group := Router.Group("pcfarm/server").Use(middleware.OperationRecord())
	query := Router.Group("pcfarm/server")
	{
		group.POST("create", pcfarmApi.CreateServerAsset)
		group.PUT("bootPolicy", pcfarmApi.UpdateBootPolicy)
		group.POST("powerAction", pcfarmApi.ExecutePowerAction)
	}
	{
		query.GET("list", pcfarmApi.GetServerAssetList)
	}
}
```

Agent routes:

```go
func (r *AgentRouter) InitAgentRouter(Router *gin.RouterGroup) {
	group := Router.Group("pcfarm/agent")
	{
		group.POST("register", pcfarmApi.Register)
		group.POST("heartbeat", pcfarmApi.Heartbeat)
	}
}
```

- [ ] **Step 5: Register group in enter files**

Modify:

```go
// server/api/v1/enter.go
PcfarmApiGroupApp pcfarm.ApiGroup
```

```go
// server/router/enter.go
PcfarmRouterGroupApp pcfarm.RouterGroup
```

Modify `server/initialize/router.go` to call:

```go
pcfarmRouter := router.RouterGroupApp.Pcfarm
pcfarmRouter.InitServerAssetRouter(PrivateGroup)
pcfarmRouter.InitAgentRouter(PublicGroup)
```

Use the exact variable style already present in `server/initialize/router.go`.

- [ ] **Step 6: Run backend route compile**

Run: `go test ./api/v1/pcfarm ./router/pcfarm ./initialize -run TestDoesNotExist`

Expected: PASS or no test files, no compile errors.

## Task 8: GORM Migration Registration

**Files:**

- Modify: `server/initialize/gorm.go` or `server/initialize/ensure_tables.go`

- [ ] **Step 1: Locate current auto-migration pattern**

Run: `rg "AutoMigrate|Ensure" server/initialize server/source`

Expected: find the existing table initialization path.

- [ ] **Step 2: Add pcfarm models to migration**

Where the project currently calls `AutoMigrate`, add:

```go
pcfarm.ServerAsset{},
pcfarm.IPPool{},
pcfarm.IPAllocation{},
pcfarm.ProvisionEvent{},
```

Import:

```go
pcfarm "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
```

- [ ] **Step 3: Run backend tests**

Run: `go test ./...`

Expected: PASS. If unrelated packages fail due to existing environment config, record exact failing package and rerun focused pcfarm packages.

## Task 9: Frontend API Wrapper

**Files:**

- Create: `web/src/api/pcfarm.js`

- [ ] **Step 1: Add API wrapper**

Create:

```js
import service from '@/utils/request'

export const createPcfarmServer = (data) => {
  return service({ url: '/pcfarm/server/create', method: 'post', data })
}

export const getPcfarmServerList = (params) => {
  return service({ url: '/pcfarm/server/list', method: 'get', params })
}

export const updatePcfarmBootPolicy = (data) => {
  return service({ url: '/pcfarm/server/bootPolicy', method: 'put', data })
}

export const executePcfarmPowerAction = (data) => {
  return service({ url: '/pcfarm/server/powerAction', method: 'post', data })
}

export const createPcfarmIPPool = (data) => {
  return service({ url: '/pcfarm/ipPool/create', method: 'post', data })
}

export const getPcfarmIPPoolList = (params) => {
  return service({ url: '/pcfarm/ipPool/list', method: 'get', params })
}
```

- [ ] **Step 2: Run frontend lint or build**

Run: `npm run build`

Expected: PASS. If the existing app has unrelated warnings, confirm there are no new import or syntax errors from `web/src/api/pcfarm.js`.

## Task 10: Server Asset Frontend Page

**Files:**

- Create: `web/src/view/pcfarm/server/index.vue`

- [ ] **Step 1: Build list page with Element Plus table**

Create a Vue SFC that imports `getPcfarmServerList`, `updatePcfarmBootPolicy`, and `executePcfarmPowerAction`.

Required state:

```js
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchInfo = ref({ keyword: '', status: '' })
```

Required columns:

- assetCode
- serialNumber
- pxeMac
- fixedIp
- powerProtocol
- bootPolicy
- status
- lastHeartbeatAt

Required actions:

- boot policy select with `local_disk`, `ubuntu_live`, `maintenance`
- power action buttons for `on`, `off`, `reboot`, `boot_pxe`

- [ ] **Step 2: Add batch action confirmation**

Use Element Plus confirm dialog before batch `off`, `reboot`, and `boot_pxe` actions:

```js
await ElMessageBox.confirm('该操作会影响选中的服务器电源状态，是否继续？', '危险操作确认', {
  confirmButtonText: '继续',
  cancelButtonText: '取消',
  type: 'warning'
})
```

- [ ] **Step 3: Verify frontend build**

Run: `npm run build`

Expected: PASS.

## Task 11: IP Pool And PXE Settings Pages

**Files:**

- Create: `web/src/view/pcfarm/ipPool/index.vue`
- Create: `web/src/view/pcfarm/pxe/index.vue`

- [ ] **Step 1: Implement IP pool page**

The page must show:

- name
- cidr
- startIp
- endIp
- gateway
- dns
- bindIface
- enabled
- allocatedCount
- availableCount

Use a dialog form for creating the single enabled pool.

- [ ] **Step 2: Implement PXE settings page**

The page must show read-only MVP settings:

- Ubuntu Live 镜像路径
- TFTP 根目录
- HTTP/NFS 地址
- dnsmasq 服务状态

Add buttons:

- 刷新配置
- 校验配置

If backend endpoints are not implemented in the same task, disable buttons and show “后端接口待接入” as button tooltip text in the UI code.

- [ ] **Step 3: Verify frontend build**

Run: `npm run build`

Expected: PASS.

## Task 12: Final Verification

**Files:**

- Review all files touched by Tasks 1-11.

- [ ] **Step 1: Format Go code**

Run: `gofmt -w server/model/pcfarm server/service/pcfarm server/api/v1/pcfarm server/router/pcfarm`

Expected: no output.

- [ ] **Step 2: Run backend tests**

Run: `go test ./...`

Expected: PASS, or document unrelated existing failures with package names and error text.

- [ ] **Step 3: Run frontend build**

Run: `npm run build`

Expected: PASS.

- [ ] **Step 4: Manual smoke test**

Start backend and frontend using existing project commands from `README.md`.

Validate:

- Create an IP pool.
- Create a server asset with PXE MAC.
- Confirm a fixed IP is assigned.
- Change boot policy to `ubuntu_live`.
- Execute a mock power action.
- Call `/pcfarm/agent/register` with matching serial number or PXE MAC.
- Confirm server status becomes `online`.

## Self-Review

- Spec coverage: covered server assets, fixed IP allocation, boot policy, PXE provider, IPMI/Redfish provider boundary, Live Agent registration, frontend pages, and verification.
- Intentional MVP gap: real `dnsmasq` file writes and real IPMI/Redfish command execution remain provider internals after the stubbed provider compiles. Implement them immediately after this plan if hardware and service paths are available.
- Placeholder scan: no unresolved placeholder tokens are used.
- Type consistency: request, model, service, API, and frontend field names use `assetCode`, `serialNumber`, `pxeMac`, `fixedIp`, `powerProtocol`, `bootPolicy`, and `status`.
