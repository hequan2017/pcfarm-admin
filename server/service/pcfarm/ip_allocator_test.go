package pcfarm

import (
	"errors"
	"testing"

	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

func TestNextAvailableIPSkipsAllocatedIP(t *testing.T) {
	pool := pcfarmModel.IPPool{StartIP: "192.168.10.10", EndIP: "192.168.10.12"}
	allocated := map[string]struct{}{
		"192.168.10.10": {},
		"192.168.10.11": {},
	}

	got, err := nextAvailableIP(pool, allocated)
	if err != nil {
		t.Fatalf("nextAvailableIP() error = %v", err)
	}
	if got != "192.168.10.12" {
		t.Fatalf("nextAvailableIP() = %q, want %q", got, "192.168.10.12")
	}
}

func TestNextAvailableIPReturnsExhaustedWhenPoolIsFull(t *testing.T) {
	pool := pcfarmModel.IPPool{StartIP: "192.168.10.10", EndIP: "192.168.10.11"}
	allocated := map[string]struct{}{
		"192.168.10.10": {},
		"192.168.10.11": {},
	}

	_, err := nextAvailableIP(pool, allocated)
	if !errors.Is(err, ErrIPPoolExhausted) {
		t.Fatalf("nextAvailableIP() error = %v, want ErrIPPoolExhausted", err)
	}
}

func TestNextAvailableIPRejectsInvalidRange(t *testing.T) {
	tests := []struct {
		name string
		pool pcfarmModel.IPPool
	}{
		{name: "reversed range", pool: pcfarmModel.IPPool{StartIP: "192.168.10.12", EndIP: "192.168.10.10"}},
		{name: "bad start ip", pool: pcfarmModel.IPPool{StartIP: "bad-ip", EndIP: "192.168.10.10"}},
		{name: "bad end ip", pool: pcfarmModel.IPPool{StartIP: "192.168.10.10", EndIP: "bad-ip"}},
		{name: "ipv6 start", pool: pcfarmModel.IPPool{StartIP: "2001:db8::1", EndIP: "192.168.10.10"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := nextAvailableIP(tt.pool, nil); err == nil {
				t.Fatal("nextAvailableIP() error = nil, want error")
			}
		})
	}
}
