package pcfarm

import (
	"errors"
	"testing"

	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

func TestPowerProviderForSelectsIPMI(t *testing.T) {
	provider, err := powerProviderFor(pcfarmModel.PowerProtocolIPMI)
	if err != nil {
		t.Fatalf("powerProviderFor() error = %v", err)
	}
	if _, ok := provider.(IPMIPowerProvider); !ok {
		t.Fatalf("powerProviderFor() = %T, want IPMIPowerProvider", provider)
	}
}

func TestPowerProviderForSelectsRedfish(t *testing.T) {
	provider, err := powerProviderFor(pcfarmModel.PowerProtocolRedfish)
	if err != nil {
		t.Fatalf("powerProviderFor() error = %v", err)
	}
	if _, ok := provider.(RedfishPowerProvider); !ok {
		t.Fatalf("powerProviderFor() = %T, want RedfishPowerProvider", provider)
	}
}

func TestPowerProviderForRejectsUnsupportedProtocol(t *testing.T) {
	_, err := powerProviderFor(pcfarmModel.PowerProtocol("wol"))
	if !errors.Is(err, ErrUnsupportedPowerProtocol) {
		t.Fatalf("powerProviderFor() error = %v, want ErrUnsupportedPowerProtocol", err)
	}
}
