package pcfarm_test

import (
	"reflect"
	"testing"
	"time"

	commonRequest "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
	"github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm/response"
)

func TestPCFarmDomainModelsExposeExpectedContracts(t *testing.T) {
	now := time.Now()
	asset := pcfarm.ServerAsset{
		AssetCode:         "srv-001",
		SerialNumber:      "sn-001",
		PxeMac:            "00:11:22:33:44:55",
		FixedIP:           "192.168.1.10",
		BmcAddress:        "192.168.1.11",
		BmcUsername:       "admin",
		BmcPasswordCipher: "cipher",
		PowerProtocol:     pcfarm.PowerProtocolIPMI,
		BootPolicy:        pcfarm.BootPolicyLocalDisk,
		Status:            pcfarm.ServerStatusOffline,
		AgentVersion:      "1.0.0",
		HardwareSummary:   "cpu=1",
		LastHeartbeatAt:   &now,
	}

	if asset.TableName() != "pcfarm_server_assets" {
		t.Fatalf("unexpected server asset table name: %s", asset.TableName())
	}

	_ = []pcfarm.BootPolicy{
		pcfarm.BootPolicyLocalDisk,
		pcfarm.BootPolicyUbuntuLive,
		pcfarm.BootPolicyMaintenance,
	}
	_ = []pcfarm.PowerProtocol{
		pcfarm.PowerProtocolIPMI,
		pcfarm.PowerProtocolRedfish,
	}
	_ = []pcfarm.ServerStatus{
		pcfarm.ServerStatusOffline,
		pcfarm.ServerStatusPXEReady,
		pcfarm.ServerStatusBooting,
		pcfarm.ServerStatusOnline,
		pcfarm.ServerStatusHeartbeatLost,
		pcfarm.ServerStatusPowerFailed,
	}

	pool := pcfarm.IPPool{Name: "default", CIDR: "192.168.1.0/24"}
	allocation := pcfarm.IPAllocation{ServerAssetID: 1, PxeMac: asset.PxeMac, IP: asset.FixedIP, Status: pcfarm.IPAllocationStatusActive}
	event := pcfarm.ProvisionEvent{ServerAssetID: 1, Type: pcfarm.ProvisionEventAssetCreated, Success: true}

	if pool.TableName() != "pcfarm_ip_pools" {
		t.Fatalf("unexpected ip pool table name: %s", pool.TableName())
	}
	if allocation.TableName() != "pcfarm_ip_allocations" {
		t.Fatalf("unexpected ip allocation table name: %s", allocation.TableName())
	}
	if event.TableName() != "pcfarm_provision_events" {
		t.Fatalf("unexpected provision event table name: %s", event.TableName())
	}

	_ = []pcfarm.IPAllocationStatus{
		pcfarm.IPAllocationStatusActive,
		pcfarm.IPAllocationStatusReleased,
	}
	_ = []pcfarm.ProvisionEventType{
		pcfarm.ProvisionEventAssetCreated,
		pcfarm.ProvisionEventIPAllocated,
		pcfarm.ProvisionEventPXERefreshed,
		pcfarm.ProvisionEventPowerAction,
		pcfarm.ProvisionEventAgentRegister,
		pcfarm.ProvisionEventHeartbeatLost,
	}

	_ = request.ServerAssetSearch{PageInfo: commonRequest.PageInfo{}, Status: string(pcfarm.ServerStatusOnline)}
	_ = request.CreateServerAsset{AssetCode: asset.AssetCode, PxeMac: asset.PxeMac, BmcPassword: "plain", PowerProtocol: string(pcfarm.PowerProtocolIPMI)}
	_ = request.UpdateBootPolicy{ID: 1, BootPolicy: string(pcfarm.BootPolicyUbuntuLive)}
	_ = request.PowerActionRequest{ID: 1, Action: "reboot"}
	_ = request.IPPoolSearch{PageInfo: commonRequest.PageInfo{}}
	_ = request.CreateIPPool{Name: pool.Name, CIDR: pool.CIDR}
	_ = request.AgentRegisterRequest{SerialNumber: asset.SerialNumber, PxeMac: asset.PxeMac, IP: asset.FixedIP, AgentVersion: asset.AgentVersion, HardwareSummary: asset.HardwareSummary, Token: "token"}
	_ = request.AgentHeartbeatRequest{SerialNumber: asset.SerialNumber, PxeMac: asset.PxeMac, Token: "token"}

	_ = response.ServerAssetResponse{Server: asset}
	_ = response.ServerAssetListItem{ServerAsset: asset, HasBmcPassword: true}
	_ = response.IPPoolSummary{Pool: pool, AllocatedCount: 1, AvailableCount: 2}
}

func TestPCFarmRequestResponseStructFieldsMatchPlan(t *testing.T) {
	assertFieldType(t, reflect.TypeOf(request.ServerAssetSearch{}), "Status", reflect.TypeOf(""))
	assertFields(t, reflect.TypeOf(request.CreateServerAsset{}), []string{
		"AssetCode",
		"SerialNumber",
		"PxeMac",
		"BmcAddress",
		"BmcUsername",
		"BmcPassword",
		"PowerProtocol",
	})
	assertFieldType(t, reflect.TypeOf(request.UpdateBootPolicy{}), "BootPolicy", reflect.TypeOf(""))
	assertFields(t, reflect.TypeOf(request.AgentRegisterRequest{}), []string{
		"SerialNumber",
		"PxeMac",
		"IP",
		"AgentVersion",
		"HardwareSummary",
		"Token",
	})
	assertFields(t, reflect.TypeOf(request.AgentHeartbeatRequest{}), []string{
		"SerialNumber",
		"PxeMac",
		"Token",
	})

	serverResponse := reflect.TypeOf(response.ServerAssetResponse{})
	assertFieldType(t, serverResponse, "Server", reflect.TypeOf(pcfarm.ServerAsset{}))
	assertJSONTag(t, serverResponse, "Server", "server")

	listItem := reflect.TypeOf(response.ServerAssetListItem{})
	assertFieldType(t, listItem, "ServerAsset", reflect.TypeOf(pcfarm.ServerAsset{}))
	assertFieldType(t, listItem, "HasBmcPassword", reflect.TypeOf(false))
	assertJSONTag(t, listItem, "HasBmcPassword", "hasBmcPassword")

	assertFields(t, reflect.TypeOf(response.IPPoolSummary{}), []string{
		"Pool",
		"AllocatedCount",
		"AvailableCount",
	})
}

func assertFields(t *testing.T, typ reflect.Type, want []string) {
	t.Helper()
	if typ.NumField() != len(want) {
		t.Fatalf("%s field count = %d, want %d", typ.Name(), typ.NumField(), len(want))
	}
	for i, name := range want {
		if typ.Field(i).Name != name {
			t.Fatalf("%s field %d = %s, want %s", typ.Name(), i, typ.Field(i).Name, name)
		}
	}
}

func assertFieldType(t *testing.T, typ reflect.Type, name string, want reflect.Type) {
	t.Helper()
	field, ok := typ.FieldByName(name)
	if !ok {
		t.Fatalf("%s missing field %s", typ.Name(), name)
	}
	if field.Type != want {
		t.Fatalf("%s.%s type = %s, want %s", typ.Name(), name, field.Type, want)
	}
}

func assertJSONTag(t *testing.T, typ reflect.Type, name string, want string) {
	t.Helper()
	field, ok := typ.FieldByName(name)
	if !ok {
		t.Fatalf("%s missing field %s", typ.Name(), name)
	}
	if got := field.Tag.Get("json"); got != want {
		t.Fatalf("%s.%s json tag = %q, want %q", typ.Name(), name, got, want)
	}
}
