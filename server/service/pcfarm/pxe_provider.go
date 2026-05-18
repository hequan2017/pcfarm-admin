package pcfarm

import (
	"fmt"

	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

type PXEProvider interface {
	Refresh(asset pcfarmModel.ServerAsset) error
	Status() (string, error)
}

type LocalDnsmasqPXEProvider struct{}

func (LocalDnsmasqPXEProvider) Refresh(asset pcfarmModel.ServerAsset) error {
	_ = renderDnsmasqHostBinding(asset)
	_ = renderBootMenu(asset.BootPolicy)
	return nil
}

func (LocalDnsmasqPXEProvider) Status() (string, error) {
	return "unknown", nil
}

func renderDnsmasqHostBinding(asset pcfarmModel.ServerAsset) string {
	return fmt.Sprintf("dhcp-host=%s,%s", asset.PxeMac, asset.FixedIP)
}

func renderBootMenu(policy pcfarmModel.BootPolicy) string {
	switch policy {
	case pcfarmModel.BootPolicyUbuntuLive:
		return "LABEL ubuntu_live\n  MENU LABEL Ubuntu Live\n  KERNEL ubuntu/vmlinuz\n  APPEND initrd=ubuntu/initrd boot=casper ip=dhcp\n"
	case pcfarmModel.BootPolicyMaintenance:
		return "LABEL maintenance\n  MENU LABEL Maintenance\n  LOCALBOOT -1\n"
	default:
		return "LABEL local_disk\n  MENU LABEL Local Disk\n  LOCALBOOT 0\n"
	}
}
