package pcfarm

import (
	"context"
	"strconv"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderPcfarmAccess = systemService.InitOrderInternal + 1

type initPcfarmAccess struct{}

func init() {
	systemService.RegisterInit(initOrderPcfarmAccess, &initPcfarmAccess{})
}

func (i *initPcfarmAccess) InitializerName() string {
	return "pcfarm_access"
}

func (i *initPcfarmAccess) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initPcfarmAccess) TableCreated(ctx context.Context) bool {
	return true
}

func (i *initPcfarmAccess) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, systemService.ErrMissingDBContext
	}
	return ctx, SyncPcfarmAccess(db)
}

func (i *initPcfarmAccess) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return hasPcfarmAccess(db)
}

func PcfarmMenus() []system.SysBaseMenu {
	return []system.SysBaseMenu{
		{
			ParentId:  0,
			Path:      "pcfarm",
			Name:      "Pcfarm",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      2,
			Meta:      system.Meta{Title: "pcfarm-admin", Icon: "server"},
		},
		{
			Path:      "server",
			Name:      "PcfarmServer",
			Hidden:    false,
			Component: "view/pcfarm/server/index.vue",
			Sort:      1,
			Meta:      system.Meta{Title: "服务器资产", Icon: "monitor"},
		},
		{
			Path:      "ipPool",
			Name:      "PcfarmIPPool",
			Hidden:    false,
			Component: "view/pcfarm/ipPool/index.vue",
			Sort:      2,
			Meta:      system.Meta{Title: "IP地址池", Icon: "connection"},
		},
		{
			Path:      "pxe",
			Name:      "PcfarmPXE",
			Hidden:    false,
			Component: "view/pcfarm/pxe/index.vue",
			Sort:      3,
			Meta:      system.Meta{Title: "PXE设置", Icon: "cpu"},
		},
		{
			Path:      "serverDetail/:id",
			Name:      "PcfarmServerDetail",
			Hidden:    true,
			Component: "view/pcfarm/server/detail.vue",
			Sort:      0,
			Meta:      system.Meta{Title: "服务器详情", Icon: "monitor", ActiveName: "PcfarmServer"},
		},
	}
}

func PcfarmApis() []system.SysApi {
	return []system.SysApi{
		{Path: "/pcfarm/server/create", Description: "创建服务器资产", ApiGroup: "服务器资产", Method: "POST"},
		{Path: "/pcfarm/server/list", Description: "获取服务器资产列表", ApiGroup: "服务器资产", Method: "GET"},
		{Path: "/pcfarm/server/bootPolicy", Description: "更新启动策略", ApiGroup: "服务器资产", Method: "PUT"},
		{Path: "/pcfarm/server/powerAction", Description: "远程电源操作", ApiGroup: "服务器资产", Method: "POST"},
		{Path: "/pcfarm/ipPool/create", Description: "创建IP地址池", ApiGroup: "IP地址池", Method: "POST"},
		{Path: "/pcfarm/ipPool/list", Description: "获取IP地址池列表", ApiGroup: "IP地址池", Method: "GET"},
		{Path: "/pcfarm/pxe/refresh", Description: "刷新PXE配置", ApiGroup: "PXE设置", Method: "POST"},
		{Path: "/pcfarm/pxe/status", Description: "获取PXE状态", ApiGroup: "PXE设置", Method: "GET"},
		{Path: "/pcfarm/agent/register", Description: "Agent注册", ApiGroup: "Agent", Method: "POST"},
		{Path: "/pcfarm/agent/heartbeat", Description: "Agent心跳", ApiGroup: "Agent", Method: "POST"},
	}
}

func SyncPcfarmAccess(db *gorm.DB) error {
	if db == nil {
		return systemService.ErrMissingDBContext
	}

	return db.Transaction(func(tx *gorm.DB) error {
		menus, err := syncPcfarmMenus(tx)
		if err != nil {
			return err
		}
		if err := syncPcfarmApis(tx); err != nil {
			return err
		}
		if err := syncAdminMenuAuthority(tx, menus); err != nil {
			return err
		}
		return syncAdminCasbin(tx)
	})
}

func SyncPcfarmAccessWithGlobalDB() error {
	return SyncPcfarmAccess(global.GVA_DB)
}

func syncPcfarmMenus(tx *gorm.DB) ([]system.SysBaseMenu, error) {
	menus := PcfarmMenus()
	if len(menus) == 0 {
		return nil, nil
	}

	parent := menus[0]
	if err := tx.Where("name = ?", parent.Name).Assign(parent).FirstOrCreate(&parent).Error; err != nil {
		return nil, errors.Wrap(err, "同步pcfarm父菜单失败")
	}
	menus[0] = parent

	for i := 1; i < len(menus); i++ {
		menu := menus[i]
		menu.ParentId = parent.ID
		if err := tx.Where("name = ?", menu.Name).Assign(menu).FirstOrCreate(&menu).Error; err != nil {
			return nil, errors.Wrap(err, "同步pcfarm子菜单失败")
		}
		menus[i] = menu
	}
	return menus, nil
}

func syncPcfarmApis(tx *gorm.DB) error {
	for _, api := range PcfarmApis() {
		if err := tx.Where("path = ? AND method = ?", api.Path, api.Method).Assign(api).FirstOrCreate(&api).Error; err != nil {
			return errors.Wrap(err, "同步pcfarm API失败")
		}
	}
	return nil
}

func syncAdminMenuAuthority(tx *gorm.DB, menus []system.SysBaseMenu) error {
	for _, menu := range menus {
		if menu.ID == 0 {
			continue
		}
		record := system.SysAuthorityMenu{
			MenuId:      strconv.Itoa(int(menu.ID)),
			AuthorityId: "888",
		}
		if err := tx.Where(record).FirstOrCreate(&record).Error; err != nil {
			return errors.Wrap(err, "同步pcfarm管理员菜单权限失败")
		}
	}
	return nil
}

func syncAdminCasbin(tx *gorm.DB) error {
	for _, api := range PcfarmApis() {
		rule := adapter.CasbinRule{Ptype: "p", V0: "888", V1: api.Path, V2: api.Method}
		if err := tx.Where(rule).FirstOrCreate(&rule).Error; err != nil {
			return errors.Wrap(err, "同步pcfarm管理员API权限失败")
		}
	}
	return nil
}

func hasPcfarmAccess(db *gorm.DB) bool {
	var menu system.SysBaseMenu
	if errors.Is(db.Where("name = ?", "PcfarmServer").First(&menu).Error, gorm.ErrRecordNotFound) {
		return false
	}

	var api system.SysApi
	if errors.Is(db.Where("path = ? AND method = ?", "/pcfarm/server/list", "GET").First(&api).Error, gorm.ErrRecordNotFound) {
		return false
	}

	var authMenu system.SysAuthorityMenu
	return !errors.Is(db.Where("sys_base_menu_id = ? AND sys_authority_authority_id = ?", strconv.Itoa(int(menu.ID)), "888").First(&authMenu).Error, gorm.ErrRecordNotFound)
}
