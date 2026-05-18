package pcfarm

import (
	"errors"
	"fmt"

	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

type PowerAction string

const (
	PowerActionOn      PowerAction = "on"
	PowerActionOff     PowerAction = "off"
	PowerActionReboot  PowerAction = "reboot"
	PowerActionBootPXE PowerAction = "boot_pxe"
)

var ErrUnsupportedPowerProtocol = errors.New("unsupported power protocol")

type PowerProvider interface {
	Execute(asset pcfarmModel.ServerAsset, action PowerAction) error
}

type IPMIPowerProvider struct{}

func (IPMIPowerProvider) Execute(asset pcfarmModel.ServerAsset, action PowerAction) error {
	return nil
}

type RedfishPowerProvider struct{}

func (RedfishPowerProvider) Execute(asset pcfarmModel.ServerAsset, action PowerAction) error {
	return nil
}

func powerProviderFor(protocol pcfarmModel.PowerProtocol) (PowerProvider, error) {
	switch protocol {
	case pcfarmModel.PowerProtocolIPMI:
		return IPMIPowerProvider{}, nil
	case pcfarmModel.PowerProtocolRedfish:
		return RedfishPowerProvider{}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedPowerProtocol, protocol)
	}
}
