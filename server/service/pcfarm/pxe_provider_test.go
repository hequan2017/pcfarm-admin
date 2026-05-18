package pcfarm

import (
	"strings"
	"testing"

	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

func TestRenderDnsmasqHostBindingIncludesMACAndIP(t *testing.T) {
	asset := pcfarmModel.ServerAsset{PxeMac: "52:54:00:12:34:56", FixedIP: "192.168.10.20"}

	got := renderDnsmasqHostBinding(asset)

	if !strings.Contains(got, asset.PxeMac) {
		t.Fatalf("renderDnsmasqHostBinding() = %q, want MAC %q", got, asset.PxeMac)
	}
	if !strings.Contains(got, asset.FixedIP) {
		t.Fatalf("renderDnsmasqHostBinding() = %q, want IP %q", got, asset.FixedIP)
	}
}

func TestRenderBootMenuUbuntuLiveContainsTitle(t *testing.T) {
	got := renderBootMenu(pcfarmModel.BootPolicyUbuntuLive)

	if !strings.Contains(got, "Ubuntu Live") {
		t.Fatalf("renderBootMenu() = %q, want Ubuntu Live", got)
	}
}
