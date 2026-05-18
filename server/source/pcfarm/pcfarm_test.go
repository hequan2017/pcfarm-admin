package pcfarm

import (
	"testing"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestPcfarmMenusExposeFrontendRoutes(t *testing.T) {
	menus := PcfarmMenus()

	if len(menus) != 5 {
		t.Fatalf("PcfarmMenus() length = %d, want 5", len(menus))
	}

	assertMenu(t, menus[0], "pcfarm", "Pcfarm", "view/routerHolder.vue", "pcfarm-admin")
	assertMenu(t, menus[1], "server", "PcfarmServer", "view/pcfarm/server/index.vue", "服务器资产")
	assertMenu(t, menus[2], "ipPool", "PcfarmIPPool", "view/pcfarm/ipPool/index.vue", "IP地址池")
	assertMenu(t, menus[3], "pxe", "PcfarmPXE", "view/pcfarm/pxe/index.vue", "PXE设置")
	assertMenu(t, menus[4], "serverDetail/:id", "PcfarmServerDetail", "view/pcfarm/server/detail.vue", "服务器详情")

	if !menus[4].Hidden {
		t.Fatalf("detail menu should be hidden")
	}
	if menus[4].Meta.ActiveName != "PcfarmServer" {
		t.Fatalf("detail active menu = %q, want PcfarmServer", menus[4].Meta.ActiveName)
	}
}

func TestPcfarmApisCoverBackendRoutes(t *testing.T) {
	apis := PcfarmApis()

	want := map[string]string{
		"POST /pcfarm/server/create":      "服务器资产",
		"GET /pcfarm/server/list":         "服务器资产",
		"PUT /pcfarm/server/bootPolicy":   "服务器资产",
		"POST /pcfarm/server/powerAction": "服务器资产",
		"POST /pcfarm/ipPool/create":      "IP地址池",
		"GET /pcfarm/ipPool/list":         "IP地址池",
		"POST /pcfarm/pxe/refresh":        "PXE设置",
		"GET /pcfarm/pxe/status":          "PXE设置",
		"POST /pcfarm/agent/register":     "Agent",
		"POST /pcfarm/agent/heartbeat":    "Agent",
	}

	if len(apis) != len(want) {
		t.Fatalf("PcfarmApis() length = %d, want %d", len(apis), len(want))
	}

	for _, api := range apis {
		key := api.Method + " " + api.Path
		group, ok := want[key]
		if !ok {
			t.Fatalf("unexpected api %s", key)
		}
		if api.ApiGroup != group {
			t.Fatalf("%s group = %q, want %q", key, api.ApiGroup, group)
		}
		delete(want, key)
	}

	if len(want) > 0 {
		t.Fatalf("missing apis: %#v", want)
	}
}

func TestSyncPcfarmAccessCreatesDisplayAndAdminPermissions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&system.SysBaseMenu{}, &system.SysApi{}, &system.SysAuthorityMenu{}, &adapter.CasbinRule{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	if err := SyncPcfarmAccess(db); err != nil {
		t.Fatalf("sync pcfarm access: %v", err)
	}

	var parent system.SysBaseMenu
	if err := db.Where("name = ?", "Pcfarm").First(&parent).Error; err != nil {
		t.Fatalf("find parent menu: %v", err)
	}

	var serverMenu system.SysBaseMenu
	if err := db.Where("name = ?", "PcfarmServer").First(&serverMenu).Error; err != nil {
		t.Fatalf("find server menu: %v", err)
	}
	if serverMenu.ParentId != parent.ID {
		t.Fatalf("server menu parent = %d, want %d", serverMenu.ParentId, parent.ID)
	}

	var api system.SysApi
	if err := db.Where("path = ? AND method = ?", "/pcfarm/server/list", "GET").First(&api).Error; err != nil {
		t.Fatalf("find list api: %v", err)
	}

	var authorityMenu system.SysAuthorityMenu
	if err := db.Where("sys_base_menu_id = ? AND sys_authority_authority_id = ?", serverMenu.ID, "888").First(&authorityMenu).Error; err != nil {
		t.Fatalf("find admin menu authority: %v", err)
	}

	var rule adapter.CasbinRule
	if err := db.Where(adapter.CasbinRule{Ptype: "p", V0: "888", V1: "/pcfarm/server/list", V2: "GET"}).First(&rule).Error; err != nil {
		t.Fatalf("find admin casbin rule: %v", err)
	}
}

func assertMenu(t *testing.T, menu system.SysBaseMenu, path, name, component, title string) {
	t.Helper()
	if menu.Path != path || menu.Name != name || menu.Component != component || menu.Meta.Title != title {
		t.Fatalf("menu mismatch: got path=%q name=%q component=%q title=%q", menu.Path, menu.Name, menu.Component, menu.Meta.Title)
	}
}
